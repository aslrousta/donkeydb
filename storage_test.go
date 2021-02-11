package donkeydb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aslrousta/donkeydb/paging"
)

func TestCreateStorage(t *testing.T) {
	rws := &paging.Memory{}
	_, err := createStorage(rws)
	assert.NoError(t, err)
}

func TestOpenStorage(t *testing.T) {
	rws := &paging.Memory{}
	createStorage(rws)
	_, err := openStorage(rws)
	assert.NoError(t, err)
}
