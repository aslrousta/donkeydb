package donkeydb

import "github.com/aslrousta/donkeydb/paging"

const (
	hashBucketBytes = 3
	hashMaxBuckets  = (pageSize - pageHeaderBytes) / hashBucketBytes

	hashFreeListOffset = 0
	hashFreeListBytes  = 3
)

type hashTable paging.Page

func (t *hashTable) Bucket(index int) int {
	return pageReadInt((*paging.Page)(t), bucketOffset(index), hashBucketBytes)
}

func (t *hashTable) SetBucket(index, value int) {
	pageWriteInt((*paging.Page)(t), bucketOffset(index), hashBucketBytes, value)
}

func (t *hashTable) FreeList() int {
	return pageReadInt((*paging.Page)(t), hashFreeListOffset, hashFreeListBytes)
}

func (t *hashTable) SetFreeList(n int) {
	pageWriteInt((*paging.Page)(t), hashFreeListOffset, hashFreeListBytes, n)
}

func bucketOffset(index int) int {
	if index >= hashMaxBuckets {
		panic("donkey: invalid bucket index")
	}
	return pageHeaderBytes + hashBucketBytes*index
}
