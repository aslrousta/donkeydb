package donkeydb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aslrousta/donkeydb/paging"
)

func TestCreateStorage(t *testing.T) {
	f := &paging.Memory{}
	_, err := createStorage(f)
	assert.NoError(t, err)
}

func TestOpenStorage(t *testing.T) {
	f := &paging.Memory{}
	createStorage(f)
	_, err := openStorage(f)
	assert.NoError(t, err)
}

func TestStorage_Get(t *testing.T) {
	t.Run("Nothing", func(t *testing.T) {
		s, _ := createStorage(&paging.Memory{})
		_, err := s.Get("key")
		assert.Equal(t, ErrNothing, err)
	})
	t.Run("Something", func(t *testing.T) {
		s, _ := createStorage(&paging.Memory{})
		s.Set("key", "value")
		v, err := s.Get("key")
		if assert.NoError(t, err) {
			value, ok := v.(string)
			if assert.True(t, ok) {
				assert.Equal(t, "value", value)
			}
		}
	})
}

func TestStorage_Set(t *testing.T) {
	s, _ := createStorage(&paging.Memory{})
	if err := s.Set("key", "value"); assert.NoError(t, err) {
		v, _ := s.Get("key")
		value, ok := v.(string)
		if assert.True(t, ok) {
			assert.Equal(t, "value", value)
		}
	}
}

func TestStorage_Del(t *testing.T) {
	t.Run("Missing", func(t *testing.T) {
		s, _ := createStorage(&paging.Memory{})
		assert.NoError(t, s.Del("key"))
	})
	t.Run("Existing", func(t *testing.T) {
		s, _ := createStorage(&paging.Memory{})
		s.Set("key", "value")
		if assert.NoError(t, s.Del("key")) {
			_, err := s.Get("key")
			assert.Equal(t, ErrNothing, err)
		}
	})
}
