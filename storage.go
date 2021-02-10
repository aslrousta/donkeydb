package donkeydb

import (
	"io"

	"github.com/aslrousta/donkeydb/paging"
)

const (
	pageSize       = 4 * 1024
	pageHeaderSize = 16
)

func createStorage(rws io.ReadWriteSeeker) (*storage, error) {
	f, err := paging.New(rws, pageSize)
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

func openStorage(rws io.ReadWriteSeeker) (*storage, error) {
	f, err := paging.New(rws, pageSize)
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
	panic("implement me")
}

func (s *storage) Set(key string, value interface{}) error {
	panic("implement me")
}
