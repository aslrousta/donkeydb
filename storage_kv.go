package donkeydb

import (
	"encoding/binary"

	"github.com/aslrousta/donkeydb/paging"
)

const (
	kvFreeOffset = 0
	kvFreeBytes  = 2
	kvNextOffset = 2
	kvNextBytes  = 3
)

type kvTable paging.Page

func (t *kvTable) Find(key string) (interface{}, bool) {
	start := t.locate(key)
	if start < 0 {
		return nil, false
	}
	page := (*paging.Page)(t)
	h := kvHeader(pageReadInt(page, start, kvHeaderBytes))
	valueStart := start + kvHeaderBytes + h.Key()
	valueEnd := valueStart + h.Value()
	value := page.Data[valueStart:valueEnd]
	switch h.Type() {
	case kvTypeString:
		return string(value), true
	case kvTypeInteger:
		n, _ := binary.Varint(value)
		return n, true
	default:
		return nil, true
	}
}

func (t *kvTable) Store(key string, value interface{}) bool {
	var data []byte
	var dataType int
	switch v := value.(type) {
	case string:
		data = []byte(v)
		dataType = kvTypeString
	case int64:
		data = make([]byte, 10)
		n := binary.PutVarint(data, v)
		data = data[:n]
		dataType = kvTypeInteger
	default:
		panic("donkey: unsupported value type")
	}
	bytes := kvHeaderBytes + len(key) + len(data)
	free := t.Free()
	if bytes > free {
		return false
	}
	offset := pageSize - free
	page := (*paging.Page)(t)
	h := kvh(dataType, len(key), len(data))
	pageWriteInt(page, offset, kvHeaderBytes, int(h))
	offset += kvHeaderBytes
	copy(page.Data[offset:], key)
	offset += len(key)
	copy(page.Data[offset:], data)
	t.SetFree(free - bytes)
	return true
}

func (t *kvTable) Del(key string) bool {
	start := t.locate(key)
	if start < 0 {
		return false
	}
	page := (*paging.Page)(t)
	h := kvHeader(pageReadInt(page, start, kvHeaderBytes))
	bytes := kvHeaderBytes + h.Key() + h.Value()
	end := start + bytes
	free := t.Free()
	copy(page.Data[start:], page.Data[end:pageSize-free])
	free += bytes
	pageWriteInt(page, pageSize-free, kvHeaderBytes, 0)
	t.SetFree(free)
	return true
}

func (t *kvTable) Free() int {
	return pageReadInt((*paging.Page)(t), kvFreeOffset, kvFreeBytes)
}

func (t *kvTable) SetFree(n int) {
	pageWriteInt((*paging.Page)(t), kvFreeOffset, kvFreeBytes, n)
}

func (t *kvTable) Next() int {
	return pageReadInt((*paging.Page)(t), kvNextOffset, kvNextBytes)
}

func (t *kvTable) SetNext(n int) {
	pageWriteInt((*paging.Page)(t), kvNextOffset, kvNextBytes, n)
}

func (t *kvTable) IsEmpty() bool {
	return t.Free() == pageSize-pageHeaderBytes
}

func (t *kvTable) locate(key string) (offset int) {
	page := (*paging.Page)(t)
	for offset = pageHeaderBytes; offset < pageSize; {
		h := kvHeader(pageReadInt(page, offset, kvHeaderBytes))
		if h.Type() == kvTypeSentinel {
			break
		}
		if h.Key() == len(key) {
			found := true
			for i := 0; i < h.Key(); i++ {
				if page.Data[offset+kvHeaderBytes+i] != key[i] {
					found = false
					break
				}
			}
			if found {
				return offset
			}
		}
		offset += kvHeaderBytes + h.Key() + h.Value()
	}
	return -1
}
