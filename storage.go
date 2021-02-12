package donkeydb

import (
	"hash/fnv"
	"io"
	"math"

	"github.com/aslrousta/donkeydb/paging"
)

const (
	pageSize        = 4 * 1024
	pageHeaderBytes = 16
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
	return &storage{
		File: f,
		Root: (*hashTable)(p),
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
	return &storage{
		File: f,
		Root: (*hashTable)(p),
	}, nil
}

type storage struct {
	File *paging.File
	Root *hashTable
}

func (s *storage) Get(key string) (interface{}, error) {
	table, bucket, err := s.table(key, false)
	if err != nil {
		return nil, err
	}
	return s.find(table, bucket, key)
}

func (s *storage) Set(key string, value interface{}) error {
	table, bucket, err := s.table(key, true)
	if err != nil {
		return err
	}
	return s.store(table, bucket, key, value)
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
	if index := int64(s.Root.Bucket(bucket)); index == 0 {
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
	} else if page, err = s.File.Read(index); err != nil {
		return nil, err
	}
	return (*hashTable)(page), nil
}

func (s *storage) find(table *hashTable, bucket int, key string) (interface{}, error) {
	index := int64(table.Bucket(bucket))
	for index != 0 {
		page, err := s.File.Read(index)
		if err != nil {
			return nil, err
		}
		kv := (*kvTable)(page)
		if value, exists := kv.Find(key); exists {
			return value, nil
		}
		index = int64(kv.Next())
	}
	return nil, ErrNothing
}

func (s *storage) store(table *hashTable, bucket int, key string, value interface{}) (err error) {
	index := int64(table.Bucket(bucket))
	var page *paging.Page
	for {
		if index == 0 {
			prevPage := page
			if page, err = s.File.Alloc(); err != nil {
				return err
			}
			(*kvTable)(page).SetFree(pageSize - pageHeaderBytes)
			if prevPage == nil {
				table.SetBucket(bucket, int(page.Index))
				if err := s.File.Write((*paging.Page)(table)); err != nil {
					return err
				}
			} else {
				previous := (*kvTable)(prevPage)
				previous.SetNext(int(page.Index))
				if err := s.File.Write(prevPage); err != nil {
					return err
				}
			}
		} else if page, err = s.File.Read(index); err != nil {
			return err
		}
		kv := (*kvTable)(page)
		if stored := kv.Store(key, value.(string)); stored {
			return s.File.Write(page)
		}
		index = int64(kv.Next())
	}
}

// hash hashes a string using fnv-1a algorithm
func hash(key string) int {
	hf := fnv.New32a()
	hf.Write([]byte(key))
	return int(hf.Sum32())
}
