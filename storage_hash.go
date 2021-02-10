package donkeydb

import "github.com/aslrousta/donkeydb/paging"

const (
	bucketSize = 3
	maxBuckets = (pageSize - pageHeaderSize) / bucketSize
)

type hashTable paging.Page

func (t *hashTable) Bucket(index int) int {
	page := (*paging.Page)(t)
	offset := bucketOffset(index)
	value := 0
	for i := 0; i < 3; i++ {
		b := int(page.Data[offset+i])
		value |= b << (8 * i)
	}
	return value
}

func (t *hashTable) SetBucket(index, value int) {
	page := (*paging.Page)(t)
	offset := bucketOffset(index)
	for i := 0; i < 3; i++ {
		page.Data[offset+i] = (byte)(value & 0xff)
		value >>= 8
	}
}

func bucketOffset(index int) int {
	if index >= maxBuckets {
		panic("donkey: out-of-bound bucket index")
	}
	return pageHeaderSize + bucketSize*index
}
