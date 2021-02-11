package donkeydb

import "github.com/aslrousta/donkeydb/paging"

const (
	kvFreeOffset = 0
	kvFreeBytes  = 2
)

type kvTable paging.Page

func (t *kvTable) Find(key string) (string, bool) {
	page := (*paging.Page)(t)
	for offset := pageHeaderBytes; offset < pageSize; {
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
				start := offset + kvHeaderBytes + h.Key()
				end := start + h.Value()
				return string(page.Data[start:end]), true
			}
		}
		offset += kvHeaderBytes + h.Key() + h.Value()
	}
	return "", false
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

func (t *kvTable) Free() int {
	return pageReadInt((*paging.Page)(t), kvFreeOffset, kvFreeBytes)
}

func (t *kvTable) SetFree(n int) {
	pageWriteInt((*paging.Page)(t), kvFreeOffset, kvFreeBytes, n)
}
