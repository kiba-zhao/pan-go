package memory

import (
	"errors"
	"slices"
	"sync"
)

type HashCode[T any] interface {
	HashCode() T
}

type BucketItem[T any] struct {
	hashcode T
}

// HashCode ...
func (bi *BucketItem[T]) HashCode() T {
	return bi.hashcode
}

// NewBucketItem ...
func NewBucketItem[T any](hash T) *BucketItem[T] {
	item := new(BucketItem[T])
	item.hashcode = hash
	return item
}

type NestBucket[H any, T any, V HashCode[T]] struct {
	*Bucket[T, V]
	*BucketItem[H]
}

// NewNestBucket ...
func NewNestBucket[H any, T any, V HashCode[T]](hash H, cmp BucketItemCompare[T]) *NestBucket[H, T, V] {
	bucket := new(NestBucket[H, T, V])
	bucket.Bucket = NewBucket[T, V](cmp)
	bucket.BucketItem = NewBucketItem[H](hash)
	return bucket
}

type Bucket[T any, V HashCode[T]] struct {
	items []V
	rw    *sync.RWMutex
	cmp   BucketItemCompare[T]
}

// GetAllHashCodes ...
func (b *Bucket[T, V]) GetHashCodes() (hashcodes []T) {
	b.rw.RLock()
	hashs := make([]T, 0)
	if len(b.items) > 0 {
		for _, item := range b.items {
			hashs = append(hashs, item.HashCode())
		}
	}

	if len(hashs) > 0 {
		hashcodes = hashs
	}

	b.rw.RUnlock()
	return
}

// Count ...
func (b *Bucket[T, V]) Count() (count int) {

	b.rw.RLock()

	count = len(b.items)

	b.rw.RUnlock()
	return
}

// GetItem ...
func (b *Bucket[T, V]) GetAll() (items []V) {

	b.rw.RLock()

	if len(b.items) > 0 {
		items = slices.Clone(b.items)
	}

	b.rw.RUnlock()
	return
}

// GetLastItem ...
func (b *Bucket[T, V]) GetLastItem() (item V) {

	b.rw.RLock()

	idx := len(b.items) - 1
	if idx >= 0 {
		item = b.items[idx]
	}

	b.rw.RUnlock()
	return
}

// GetItem ...
func (b *Bucket[T, V]) GetItem(hash T) (item V) {

	b.rw.RLock()

	idx, existed := findBucketItemIdx(b.items, hash, b.cmp)
	if existed {
		item = b.items[idx]
	}

	b.rw.RUnlock()
	return
}

// SetItem ...
func (b *Bucket[T, V]) SetItem(item V) {
	b.rw.Lock()

	hash := item.HashCode()
	idx, existed := findBucketItemIdx(b.items, hash, b.cmp)
	if existed {
		b.items[idx] = item
	} else {
		if idx < 0 {
			b.items = append(b.items, item)
		} else {
			b.items = slices.Insert(b.items, idx, item)
		}
	}
	b.rw.Unlock()
}

// GetOrAddItem ...
func (b *Bucket[T, V]) GetOrAddItem(item V) (ritem V, existed bool) {
	b.rw.Lock()

	hash := item.HashCode()
	idx, existed := findBucketItemIdx(b.items, hash, b.cmp)
	if existed {
		ritem = b.items[idx]
	} else {
		if idx < 0 {
			b.items = append(b.items, item)
		} else {
			b.items = slices.Insert(b.items, idx, item)
		}
		ritem = item
	}

	b.rw.Unlock()
	return
}

// AddItem ...
func (b *Bucket[T, V]) AddItem(item V) (err error) {
	b.rw.Lock()

	hash := item.HashCode()
	idx, existed := findBucketItemIdx(b.items, hash, b.cmp)
	if existed {
		err = errors.New("Bucket Item already existsed")
	} else {
		if idx < 0 {
			b.items = append(b.items, item)
		} else {
			b.items = slices.Insert(b.items, idx, item)
		}
	}

	b.rw.Unlock()
	return
}

// RemoveItem ...
func (b *Bucket[T, V]) RemoveItem(item V) {

	b.rw.Lock()

	lastIdx := len(b.items) - 1

	existed := lastIdx >= 0
	idx := -1
	if existed {
		hash := item.HashCode()
		idx, existed = findBucketItemIdx(b.items, hash, b.cmp)
	}
	if existed {
		items := make([]V, lastIdx)
		if idx == lastIdx && idx != 0 {
			copy(items, b.items[:lastIdx])
		} else if idx != lastIdx {
			if idx != 0 {
				copy(items, b.items[:idx])
			}
			copy(items, b.items[idx+1:])
		}
		b.items = items
	}

	b.rw.Unlock()
}

type BucketItemCompare[T any] func(prev, next T) int

// findBucketItemIdx ...
func findBucketItemIdx[T any, V HashCode[T]](items []V, hash T, compare BucketItemCompare[T]) (idx int, existed bool) {

	idx = -1
	maxIdx := len(items) - 1
	minIdx := 0
	prevIdx := -1
	existed = false
	for minIdx <= maxIdx {
		midIdx := (minIdx + maxIdx) / 2
		midItem := items[midIdx]
		cmp := compare(hash, midItem.HashCode())
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

// NewBucket ...
func NewBucket[T any, V HashCode[T]](cmp BucketItemCompare[T]) *Bucket[T, V] {
	bucket := new(Bucket[T, V])
	bucket.items = make([]V, 0)
	bucket.rw = new(sync.RWMutex)
	bucket.cmp = cmp
	return bucket
}
