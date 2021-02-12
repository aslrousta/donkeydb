package donkeydb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aslrousta/donkeydb"
	"github.com/aslrousta/donkeydb/paging"
)

func TestDatabase_Get(t *testing.T) {
	t.Run("Nothing", func(t *testing.T) {
		if d, err := donkeydb.Create(&paging.Memory{}); assert.NoError(t, err) {
			_, err := d.Get("key")
			assert.Equal(t, donkeydb.ErrNothing, err)
		}
	})
	t.Run("Something", func(t *testing.T) {
		if d, err := donkeydb.Create(&paging.Memory{}); assert.NoError(t, err) {
			d.Set("key", "value")
			v, err := d.Get("key")
			if assert.NoError(t, err) {
				value, ok := v.(string)
				if assert.True(t, ok) {
					assert.Equal(t, "value", value)
				}
			}
		}
	})
}

func TestDatabase_Set(t *testing.T) {
	if d, err := donkeydb.Create(&paging.Memory{}); assert.NoError(t, err) {
		if err := d.Set("key", "value"); assert.NoError(t, err) {
			v, _ := d.Get("key")
			value, ok := v.(string)
			if assert.True(t, ok) {
				assert.Equal(t, "value", value)
			}
		}
	}
}
