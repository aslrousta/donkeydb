package donkeydb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVTable(t *testing.T) {
	kv := &kvTable{Data: make([]byte, pageSize)}
	kv.SetFree(pageSize - pageHeaderBytes)
	if !assert.True(t, kv.Store("one", "first")) ||
		!assert.True(t, kv.Store("two", int64(2))) {
		return
	}
	t.Run("Existing", func(t *testing.T) {
		value, ok := kv.Find("two")
		if assert.True(t, ok) {
			assert.Equal(t, int64(2), value)
		}
	})
	t.Run("Missing", func(t *testing.T) {
		_, ok := kv.Find("three")
		assert.False(t, ok)
	})
	t.Run("Delete", func(t *testing.T) {
		if assert.True(t, kv.Del("two")) {
			_, ok := kv.Find("two")
			assert.False(t, ok)
		}
	})
}
