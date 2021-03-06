package donkeydb

import (
	"hash/fnv"
	"io"
	"math"
	"sync"

	lru "github.com/hashicorp/golang-lru"

	"github.com/aslrousta/donkeydb/paging"
)

const (
	pageSize        = 4 * 1024
	pageHeaderBytes = 16
	pageCacheSize   = 8
)

func createStorage(s io.ReadWriteSeeker) (*storage, error) {
	f, err := paging.New(s, pageSize)
	if err != nil {
		return nil, err
	}
	p, err := f.Alloc()
	if err != nil {
		return nil, err
	}
	cache, err := lru.New(pageCacheSize)
	if err != nil {
		return nil, err
	}
	return &storage{
		File:  f,
		Cache: cache,
		Root:  (*hashTable)(p),
	}, nil
}

func openStorage(s io.ReadWriteSeeker) (*storage, error) {
	f, err := paging.New(s, pageSize)
	if err != nil {
		return nil, err
	}
	p, err := f.Read(0)
	if err != nil {
		return nil, err
	}
	cache, err := lru.New(pageCacheSize)
	if err != nil {
		return nil, err
	}
	return &storage{
		File:  f,
		Cache: cache,
		Root:  (*hashTable)(p),
	}, nil
}

type storage struct {
	Mutex      sync.RWMutex
	File       *paging.File
	Root       *hashTable
	CacheMutex sync.Mutex
	Cache      *lru.Cache
}

func (s *storage) Get(key string) (interface{}, error) {
	if err := checkKeyLen(len(key)); err != nil {
		return nil, err
	}
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	table, bucket, err := s.table(key, false)
	if err != nil {
		return nil, err
	}
	return s.find(table, bucket, key)
}

func (s *storage) Set(key string, value interface{}) error {
	if err := checkKeyLen(len(key)); err != nil {
		return err
	}
	switch v := value.(type) {
	case string:
		if len(v) > kvHeaderValueMaxLen || 2*len(v) > pageSize {
			return ErrValueTooLong
		}
	case int64:
		// Let is pass.
	default:
		return ErrUnsuppValue
	}
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	table, bucket, err := s.table(key, true)
	if err != nil {
		return err
	}
	if err := s.del(table, bucket, key); err != nil {
		return err
	}
	return s.store(table, bucket, key, value)
}

func (s *storage) Del(key string) error {
	if err := checkKeyLen(len(key)); err != nil {
		return err
	}
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	table, bucket, err := s.table(key, false)
	if err != nil {
		return err
	}
	return s.del(table, bucket, key)
}

func (s *storage) table(key string, create bool) (*hashTable, int, error) {
	h := hash(key)
	interval := (math.MaxInt32 + hashMaxBuckets - 1) / hashMaxBuckets
	bucket := h / interval
	table, err := s.secondary(bucket, true)
	if err != nil {
		return nil, 0, err
	}
	hashStart := bucket * interval
	interval = (interval + hashMaxBuckets - 1) / hashMaxBuckets
	return table, (h - hashStart) / interval, nil
}

func (s *storage) secondary(bucket int, create bool) (*hashTable, error) {
	var page *paging.Page
	var err error
	if index := s.Root.Bucket(bucket); index == 0 {
		if !create {
			return nil, ErrNothing
		}
		if page, err = s.File.Alloc(); err != nil {
			return nil, err
		}
		s.Root.SetBucket(bucket, int(page.Index))
		if err := s.File.Write((*paging.Page)(s.Root)); err != nil {
			return nil, err
		}
	} else if page, err = s.load(index); err != nil {
		return nil, err
	}
	return (*hashTable)(page), nil
}

func (s *storage) find(table *hashTable, bucket int, key string) (interface{}, error) {
	index := table.Bucket(bucket)
	for index != 0 {
		page, err := s.load(index)
		if err != nil {
			return nil, err
		}
		kv := (*kvTable)(page)
		if value, exists := kv.Find(key); exists {
			return value, nil
		}
		index = kv.Next()
	}
	return nil, ErrNothing
}

func (s *storage) store(table *hashTable, bucket int, key string, value interface{}) (err error) {
	index := table.Bucket(bucket)
	var page *paging.Page
	for {
		if index == 0 {
			prevPage := page
			if page, err = s.alloc(); err != nil {
				return err
			}
			(*kvTable)(page).SetFree(pageSize - pageHeaderBytes)
			if prevPage == nil {
				table.SetBucket(bucket, int(page.Index))
				if err := s.File.Write((*paging.Page)(table)); err != nil {
					return err
				}
			} else {
				prev := (*kvTable)(prevPage)
				prev.SetNext(int(page.Index))
				if err := s.File.Write(prevPage); err != nil {
					return err
				}
			}
		} else if page, err = s.load(index); err != nil {
			return err
		}
		kv := (*kvTable)(page)
		if stored := kv.Store(key, value.(string)); stored {
			return s.File.Write(page)
		}
		index = kv.Next()
	}
}

func (s *storage) del(table *hashTable, bucket int, key string) error {
	var prev *kvTable
	index := table.Bucket(bucket)
	for index != 0 {
		page, err := s.load(index)
		if err != nil {
			return err
		}
		kv := (*kvTable)(page)
		if deleted := kv.Del(key); deleted {
			if !kv.IsEmpty() {
				return s.File.Write(page)
			}
			return s.dealloc(kv, prev, table, bucket)
		}
		prev = kv
		index = kv.Next()
	}
	return nil
}

func (s *storage) load(index int) (*paging.Page, error) {
	s.CacheMutex.Lock()
	defer s.CacheMutex.Unlock()
	if v, ok := s.Cache.Get(index); ok {
		return v.(*paging.Page), nil
	}
	page, err := s.File.Read(int64(index))
	if err != nil {
		return nil, err
	}
	s.Cache.Add(index, page)
	return page, nil
}

func (s *storage) alloc() (*paging.Page, error) {
	if s.Root.FreeList() == 0 {
		return s.File.Alloc()
	}
	page, err := s.load(s.Root.FreeList())
	if err != nil {
		return nil, err
	}
	kv := (*kvTable)(page)
	s.Root.SetFreeList(kv.Next())
	if err := s.File.Write((*paging.Page)(s.Root)); err != nil {
		return nil, err
	}
	return page, nil
}

func (s *storage) dealloc(kv, prev *kvTable, table *hashTable, bucket int) error {
	if prev == nil {
		table.SetBucket(bucket, kv.Next())
		if err := s.File.Write((*paging.Page)(table)); err != nil {
			return err
		}
	} else {
		prev.SetNext(kv.Next())
		if err := s.File.Write((*paging.Page)(prev)); err != nil {
			return err
		}
	}
	kv.SetNext(s.Root.FreeList())
	if err := s.File.Write((*paging.Page)(kv)); err != nil {
		return err
	}
	s.Root.SetFreeList(int(kv.Index))
	return s.File.Write((*paging.Page)(s.Root))
}

func checkKeyLen(keyLen int) error {
	switch {
	case keyLen == 0:
		return ErrKeyTooShort
	case keyLen > kvHeaderKeyMaxLen:
		return ErrKeyTooLong
	default:
		return nil
	}
}

// hash hashes a string using fnv-1a algorithm
func hash(key string) int {
	hf := fnv.New32a()
	hf.Write([]byte(key))
	return int(hf.Sum32())
}
