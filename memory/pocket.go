package memory

import (
	"slices"
	"sync"
)

type Pocket[T comparable] struct {
	items []T
	rw    *sync.RWMutex
}

// NewPocket ...
func NewPocket[T comparable]() *Pocket[T] {
	pocket := new(Pocket[T])
	pocket.items = make([]T, 0)
	pocket.rw = new(sync.RWMutex)
	return pocket
}

// GetAll ...
func (p *Pocket[T]) GetAll() (items []T) {

	p.rw.RLock()

	if len(p.items) > 0 {
		items = slices.Clone(p.items)
	}

	p.rw.RUnlock()
	return
}

// Add ...
func (p *Pocket[T]) Add(items ...T) {
	p.rw.Lock()
	p.items = append(p.items, items...)
	p.rw.Unlock()
}

// Remove ...
func (p *Pocket[T]) Remove(item T) {
	p.rw.Lock()
	idx := slices.Index(p.items, item)
	if idx >= 0 {
		lastIdx := len(p.items) - 1
		if lastIdx != idx {
			p.items[idx] = p.items[lastIdx]
		}
		p.items = p.items[:lastIdx]
	}

	p.rw.Unlock()
}
