package donkeydb

import (
	"errors"
	"io"
)

var (
	// ErrNothing reports a non-existing key.
	ErrNothing = errors.New("donkey: key not found")
)

// Create creates a new key-value database.
func Create(s io.ReadWriteSeeker) (*Database, error) {
	storage, err := createStorage(s)
	if err != nil {
		return nil, err
	}
	return &Database{storage: storage}, nil
}

// Open opens an existing key-value database.
func Open(s io.ReadWriteSeeker) (*Database, error) {
	storage, err := openStorage(s)
	if err != nil {
		return nil, err
	}
	return &Database{storage: storage}, nil
}

// Database is a disk-backed key-value database.
type Database struct {
	storage *storage
}

// Get retrieves a value for a given key.
func (d *Database) Get(key string) (interface{}, error) {
	return d.storage.Get(key)
}

// Set stores a value for a given key.
func (d *Database) Set(key string, value interface{}) error {
	return d.storage.Set(key, value)
}

// Del removes a value for a given key.
func (d *Database) Del(key string) error {
	return d.storage.Del(key)
}
