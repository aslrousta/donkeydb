package donkeydb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVHeader(t *testing.T) {
	v := kvh(1, 2, 3)
	assert.Equal(t, 1, v.Type())
	assert.Equal(t, 2, v.Key())
	assert.Equal(t, 3, v.Value())
}
