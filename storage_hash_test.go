package donkeydb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashTable(t *testing.T) {
	h := &hashTable{Data: make([]byte, pageSize)}
	magic := (1 << 16) | (2 << 8) | 3
	h.SetBucket(0, magic)
	assert.Equal(t, magic, h.Bucket(0))
}
