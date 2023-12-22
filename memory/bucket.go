package memory

import (
	"slices"
	"sync"
)

type BucketBlockCompare[T any] func(prev, next T) int

type BucketItem[T any, V any] struct {
	block *BucketBlock[T, V]
	idx   int
	value T
	rw    *sync.RWMutex
}

// isExpired ...
func (bi *BucketItem[T, V]) Expired() (expired bool) {

	bi.rw.RLock()
	expired = bi.idx < 0
	bi.rw.RUnlock()

	return
}

// Value ...
func (bi *BucketItem[T, V]) Value() (val T) {

	bi.rw.RLock()
	val = bi.value
	bi.rw.RUnlock()

	return
}

// Block ...
// func (bi *BucketItem[T, V]) Block() *BucketBlock[T, V] {
// 	return bi.block
// }

type BucketBlock[T any, V any] struct {
	id    V
	idx   int
	items []*BucketItem[T, V]
}

type Bucket[T any, V any] struct {
	blocks []*BucketBlock[T, V]
	rw     *sync.RWMutex
	cmp    BucketBlockCompare[V]
}

// FindBlock ...
// func (b *Bucket[T, V]) FindBlock(targetId V) (block *BucketBlock[T, V]) {

// 	b.rw.RLock()
// 	idx, existed := findBlockIdx(b.blocks, targetId, b.cmp)
// 	if existed {
// 		block = b.blocks[idx]
// 	}
// 	b.rw.RUnlock()

// 	return
// }

// FindItems ...
func (b *Bucket[T, V]) FindBlockItems(targetId V) (items []*BucketItem[T, V]) {

	b.rw.RLock()
	idx, existed := findBlockIdx(b.blocks, targetId, b.cmp)
	if existed {
		block := b.blocks[idx]
		if len(block.items) > 0 {
			items = slices.Clone(block.items)
		}
	}
	b.rw.RUnlock()

	return
}

// FindItems ...
func (b *Bucket[T, V]) FindBlockItem(targetId V) (item *BucketItem[T, V]) {

	b.rw.RLock()
	idx, existed := findBlockIdx(b.blocks, targetId, b.cmp)
	if existed {
		block := b.blocks[idx]
		for _, bitem := range block.items {
			if bitem.idx >= 0 {
				item = bitem
				break
			}
		}

	}
	b.rw.RUnlock()

	return
}

// PutItem ...
func (b *Bucket[T, V]) PutItem(targetId V, value T) (item *BucketItem[T, V]) {

	item = new(BucketItem[T, V])
	item.value = value
	item.rw = new(sync.RWMutex)

	b.rw.Lock()

	idx, existed := findBlockIdx(b.blocks, targetId, b.cmp)

	if existed {
		block := b.blocks[idx]
		item.block = block
		item.idx = len(block.items)
		block.items = append(block.items, item)
	} else {
		block := new(BucketBlock[T, V])
		item.block = block
		item.idx = 0
		block.id = targetId
		block.items = make([]*BucketItem[T, V], 0, 1)
		block.items = append(block.items, item)
		if idx < 0 {
			block.idx = len(b.blocks)
			b.blocks = append(b.blocks, block)
		} else {
			block.idx = idx
			b.blocks = slices.Insert(b.blocks, idx, block)
			for _, nblock := range b.blocks[idx+1:] {
				nblock.idx++
			}
		}
	}

	b.rw.Unlock()

	return
}

// RemoveItem ...
func (b *Bucket[T, V]) RemoveItem(item *BucketItem[T, V]) {

	b.rw.Lock()
	item.rw.Lock()

	block := item.block
	lastIdx := len(block.items) - 1
	idx := item.idx
	if item.idx >= 0 && lastIdx >= 0 {
		items := make([]*BucketItem[T, V], lastIdx)
		if idx == lastIdx {
			copy(items, block.items[:lastIdx])
		} else {
			if idx != 0 {
				copy(items, block.items[:idx])
			}
			copy(items, block.items[idx+1:])
			for _, nitem := range items[idx:] {
				nitem.idx--
			}
		}
		block.items = items
	}

	item.idx = -1
	item.rw.Unlock()
	b.rw.Unlock()

	return
}

// NewBucket ...
func NewBucket[T any, V any](cmp BucketBlockCompare[V]) *Bucket[T, V] {

	bucket := new(Bucket[T, V])
	bucket.blocks = make([]*BucketBlock[T, V], 0)
	bucket.rw = new(sync.RWMutex)
	bucket.cmp = cmp
	return bucket
}

// findBlockIdx ...
func findBlockIdx[T any, V any](blocks []*BucketBlock[T, V], targetId V, compare BucketBlockCompare[V]) (idx int, existed bool) {

	idx = -1
	maxIdx := len(blocks) - 1
	minIdx := 0
	prevIdx := -1
	existed = false
	for minIdx <= maxIdx {
		midIdx := (minIdx + maxIdx) / 2
		block := blocks[midIdx]
		cmp := compare(block.id, targetId)
		if cmp == 0 {
			existed = true
			idx = midIdx
			break
		}

		if cmp < 0 {
			prevIdx = midIdx
			minIdx = midIdx + 1
			continue
		}
		if prevIdx+1 == midIdx {
			idx = midIdx
			break
		}

		maxIdx = midIdx - 1

	}
	return
}
