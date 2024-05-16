package cache

import (
	"errors"
	"slices"
	"sync"
)

type HashCode[T any] interface {
	HashCode() T
}

type Bucket[T any, V HashCode[T]] interface {
	HashCodes() []T
	Size() int
	Swap(item V) (previous V, ok bool)
	Search(hash T) (item V, ok bool)
	SearchOrStore(item V) (ritem V, ok bool)
	Store(item V) (err error)
	Delete(item V)
	At(index int) (item V, ok bool)
	Items() []V
}

type NestBucket[H any, T any, V HashCode[T]] struct {
	Bucket[T, V]
	Code H
}

func (b *NestBucket[H, T, V]) HashCode() H {
	return b.Code
}

type SimpleBucket[T any, V HashCode[T]] struct {
	items []V
	cmp   BucketItemCompare[T]
}

// NewBucket ...
func NewBucket[T any, V HashCode[T]](cmp BucketItemCompare[T]) Bucket[T, V] {
	bucket := &SimpleBucket[T, V]{}
	bucket.items = make([]V, 0)
	bucket.cmp = cmp
	return bucket
}

// HashCodes ...
func (b *SimpleBucket[T, V]) HashCodes() (hashcodes []T) {
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
func (b *SimpleBucket[T, V]) Size() int {
	return len(b.items)
}

// Items ...
func (b *SimpleBucket[T, V]) Items() []V {
	return slices.Clone(b.items)
}

// At ...
func (b *SimpleBucket[T, V]) At(index int) (item V, ok bool) {

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
func (b *SimpleBucket[T, V]) Search(hash T) (item V, ok bool) {

	idx, ok := SearchBucketIndex(b.items, hash, b.cmp)
	if ok {
		item = b.items[idx]
	}

	return
}

// Swap ...
func (b *SimpleBucket[T, V]) Swap(item V) (previous V, ok bool) {

	hash := item.HashCode()
	idx, ok := SearchBucketIndex(b.items, hash, b.cmp)
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
func (b *SimpleBucket[T, V]) SearchOrStore(item V) (ritem V, ok bool) {

	hash := item.HashCode()
	idx, ok := SearchBucketIndex(b.items, hash, b.cmp)
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
func (b *SimpleBucket[T, V]) Store(item V) (err error) {

	hash := item.HashCode()
	idx, ok := SearchBucketIndex(b.items, hash, b.cmp)

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
func (b *SimpleBucket[T, V]) Delete(item V) {

	lastIdx := len(b.items) - 1

	ok := lastIdx >= 0
	idx := -1
	if ok {
		hash := item.HashCode()
		idx, ok = SearchBucketIndex(b.items, hash, b.cmp)
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

type RWBucket[T any, V HashCode[T]] struct {
	bucket Bucket[T, V]
	rw     sync.RWMutex
}

func WrapSyncBucket[T any, V HashCode[T]](bucket Bucket[T, V]) Bucket[T, V] {
	bucket_ := &RWBucket[T, V]{}
	bucket_.bucket = bucket
	return bucket_
}

func (b *RWBucket[T, V]) HashCodes() []T {
	b.rw.RLock()
	defer b.rw.RUnlock()

	return b.bucket.HashCodes()
}

// Size ...
func (b *RWBucket[T, V]) Size() int {

	b.rw.RLock()
	defer b.rw.RUnlock()

	return b.bucket.Size()
}

// Items ...
func (b *RWBucket[T, V]) Items() []V {

	b.rw.RLock()
	defer b.rw.RUnlock()

	return b.bucket.Items()
}

// At ...
func (b *RWBucket[T, V]) At(index int) (V, bool) {

	b.rw.RLock()
	defer b.rw.RUnlock()

	return b.bucket.At(index)
}

// Search ...
func (b *RWBucket[T, V]) Search(hash T) (V, bool) {

	b.rw.RLock()
	defer b.rw.RUnlock()

	return b.bucket.Search(hash)
}

// Swap ...
func (b *RWBucket[T, V]) Swap(item V) (V, bool) {
	b.rw.Lock()
	defer b.rw.Unlock()

	return b.bucket.Swap(item)
}

// SearchOrStore ...
func (b *RWBucket[T, V]) SearchOrStore(item V) (V, bool) {
	b.rw.Lock()
	defer b.rw.Unlock()

	return b.bucket.SearchOrStore(item)
}

// Store ...
func (b *RWBucket[T, V]) Store(item V) error {
	b.rw.Lock()
	defer b.rw.Unlock()

	return b.bucket.Store(item)
}

// Delete ...
func (b *RWBucket[T, V]) Delete(item V) {

	b.rw.Lock()
	defer b.rw.Unlock()

	b.bucket.Delete(item)
}

type BucketItemCompare[T any] func(prev, next T) int

// SearchBucketIndex ...
func SearchBucketIndex[T any, V HashCode[T]](items []V, hash T, compare BucketItemCompare[T]) (idx int, ok bool) {

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
