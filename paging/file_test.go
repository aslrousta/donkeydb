package paging_test

import (
	"math/rand"
	"testing"

	"github.com/aslrousta/donkeydb/paging"
)

func TestFile(t *testing.T) {
	orig := make([]byte, 128)
	rand.Read(orig)

	f, err := paging.New(&paging.Memory{}, len(orig))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		page, err := f.Alloc()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		copy(page.Data, orig)
		if err := f.Write(page); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	page, err := f.Read(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < len(page.Data); i++ {
		if page.Data[i] != orig[i] {
			t.Fatalf("unexpected byte at: %d", i)
		}
	}
}
