package donkeydb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVTable(t *testing.T) {
	kv := &kvTable{Data: make([]byte, pageSize)}
	kv.SetFree(pageSize - pageHeaderBytes)
	if !assert.True(t, kv.Store("one", "first")) ||
		!assert.True(t, kv.Store("two", "second")) {
		return
	}
	t.Run("Existing", func(t *testing.T) {
		value, ok := kv.Find("two")
		if assert.True(t, ok) {
			assert.Equal(t, "second", value)
		}
	})
	t.Run("Missing", func(t *testing.T) {
		_, ok := kv.Find("three")
		assert.False(t, ok)
	})
}
