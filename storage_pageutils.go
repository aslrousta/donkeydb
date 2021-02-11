package donkeydb

import "github.com/aslrousta/donkeydb/paging"

func pageReadInt(p *paging.Page, offset, size int) int {
	value := 0
	for i := 0; i < size; i++ {
		b := int(p.Data[offset+i])
		value |= b << (8 * i)
	}
	return value
}

func pageWriteInt(p *paging.Page, offset, size, value int) {
	for i := 0; i < size; i++ {
		p.Data[offset+i] = (byte)(value & 0xff)
		value >>= 8
	}
}
