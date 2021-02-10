package donkeydb

import (
	"testing"

	"github.com/aslrousta/donkeydb/paging"
)

func TestCreateStorage(t *testing.T) {
	rws := &paging.Memory{}
	if _, err := createStorage(rws); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOpenStorage(t *testing.T) {
	rws := &paging.Memory{}
	createStorage(rws)
	if _, err := openStorage(rws); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
