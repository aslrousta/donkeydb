package donkeydb

import "errors"

var (
	// ErrNothing reports a non-existing key.
	ErrNothing = errors.New("donkey: key not found")
)

// New instantiates a new key-value database.
func New() *Database {
	return &Database{
		storage: newStorage(),
	}
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
