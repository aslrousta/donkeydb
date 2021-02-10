package donkeydb

import "testing"

func TestHashTable(t *testing.T) {
	h := &hashTable{
		Data: make([]byte, pageHeaderSize+bucketSize),
	}
	magic := (1 << 16) | (2 << 8) | 3
	h.SetBucket(0, magic)
	if h.Bucket(0) != magic {
		t.Fatal("bucket data mismatch")
	}
}
