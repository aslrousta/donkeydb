package donkeydb

import "github.com/aslrousta/donkeydb/paging"

const (
	kvFreeOffset = 0
	kvFreeBytes  = 2
	kvNextOffset = 2
	kvNextBytes  = 3
)

type kvTable paging.Page

func (t *kvTable) Find(key string) (string, bool) {
	start := t.locate(key)
	if start < 0 {
		return "", false
	}
	page := (*paging.Page)(t)
	h := kvHeader(pageReadInt(page, start, kvHeaderBytes))
	valueStart := start + kvHeaderBytes + h.Key()
	valueEnd := valueStart + h.Value()
	return string(page.Data[valueStart:valueEnd]), true
}

func (t *kvTable) Store(key, value string) bool {
	bytes := kvHeaderBytes + len(key) + len(value)
	free := t.Free()
	if bytes > free {
		return false
	}
	offset := pageSize - free
	page := (*paging.Page)(t)
	h := kvh(kvTypeString, len(key), len(value))
	pageWriteInt(page, offset, kvHeaderBytes, int(h))
	offset += kvHeaderBytes
	copy(page.Data[offset:], key)
	offset += len(key)
	copy(page.Data[offset:], value)
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
