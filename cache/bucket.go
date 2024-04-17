package cache

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
	bucket.BucketItem = NewBucketItem(hash)
	return bucket
}

type Bucket[T any, V HashCode[T]] struct {
	items []V
	rw    *sync.RWMutex
	cmp   BucketItemCompare[T]
}

// NewBucket ...
func NewBucket[T any, V HashCode[T]](cmp BucketItemCompare[T]) *Bucket[T, V] {
	bucket := new(Bucket[T, V])
	bucket.items = make([]V, 0)
	bucket.rw = new(sync.RWMutex)
	bucket.cmp = cmp
	return bucket
}

// HashCodes ...
func (b *Bucket[T, V]) HashCodes() (hashcodes []T) {
	b.rw.RLock()
	defer b.rw.RUnlock()
	hashs := make([]T, 0)
	if len(b.items) > 0 {
		for _, item := range b.items {
			hashs = append(hashs, item.HashCode())
		}
	}

	if len(hashs) > 0 {
		hashcodes = hashs
	}

	return
}

// Size ...
func (b *Bucket[T, V]) Size() (size int) {

	b.rw.RLock()
	defer b.rw.RUnlock()

	size = len(b.items)

	return
}

// GetItem ...
func (b *Bucket[T, V]) Items() (items []V) {

	b.rw.RLock()
	defer b.rw.RUnlock()

	if len(b.items) > 0 {
		items = slices.Clone(b.items)
	}

	return
}

// At ...
func (b *Bucket[T, V]) At(index int) (item V, ok bool) {

	b.rw.RLock()
	defer b.rw.RUnlock()

	idx := index
	if idx < 0 {
		idx = len(b.items) + idx
	}

	if idx >= 0 && idx < len(b.items) {
		item = b.items[idx]
		ok = true
	}

	return
}

// Search ...
func (b *Bucket[T, V]) Search(hash T) (item V, ok bool) {

	b.rw.RLock()
	defer b.rw.RUnlock()

	idx, ok := indexOfBucket(b.items, hash, b.cmp)
	if ok {
		item = b.items[idx]
	}

	return
}

// Swap ...
func (b *Bucket[T, V]) Swap(item V) (previous V, ok bool) {
	b.rw.Lock()
	defer b.rw.Unlock()

	hash := item.HashCode()
	idx, ok := indexOfBucket(b.items, hash, b.cmp)
	if ok {
		previous = b.items[idx]
		b.items[idx] = item
	} else {
		if idx < 0 {
			b.items = append(b.items, item)
		} else {
			b.items = slices.Insert(b.items, idx, item)
		}
	}
	return
}

// SearchOrStore ...
func (b *Bucket[T, V]) SearchOrStore(item V) (ritem V, ok bool) {
	b.rw.Lock()
	defer b.rw.Unlock()

	hash := item.HashCode()
	idx, ok := indexOfBucket(b.items, hash, b.cmp)
	if ok {
		ritem = b.items[idx]
	} else {
		if idx < 0 {
			b.items = append(b.items, item)
		} else {
			b.items = slices.Insert(b.items, idx, item)
		}
		ritem = item
	}

	return
}

// Store ...
func (b *Bucket[T, V]) Store(item V) (err error) {
	b.rw.Lock()
	defer b.rw.Unlock()

	hash := item.HashCode()
	idx, ok := indexOfBucket(b.items, hash, b.cmp)

	if ok {
		err = errors.New("Bucket Item already existsed")
	} else {
		if idx < 0 {
			b.items = append(b.items, item)
		} else {
			b.items = slices.Insert(b.items, idx, item)
		}
	}

	return
}

// Delete ...
func (b *Bucket[T, V]) Delete(item V) {

	b.rw.Lock()
	defer b.rw.Unlock()

	lastIdx := len(b.items) - 1

	ok := lastIdx >= 0
	idx := -1
	if ok {
		hash := item.HashCode()
		idx, ok = indexOfBucket(b.items, hash, b.cmp)
	}
	if ok {
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

}

type BucketItemCompare[T any] func(prev, next T) int

// indexOfBucket ...
func indexOfBucket[T any, V HashCode[T]](items []V, hash T, compare BucketItemCompare[T]) (idx int, ok bool) {

	idx = -1
	maxIdx := len(items) - 1
	minIdx := 0
	prevIdx := -1
	ok = false
	for minIdx <= maxIdx {
		midIdx := (minIdx + maxIdx) / 2
		midItem := items[midIdx]
		cmp := compare(midItem.HashCode(), hash)
		if cmp == 0 {
			ok = true
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
