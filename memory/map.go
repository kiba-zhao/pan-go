package memory

import "sync"

type Map[T any, V any] struct {
	internalMap *sync.Map
}

// Store ...
func (m *Map[T, V]) Store(key T, value V) {
	m.internalMap.Store(key, value)
}

// Range ...
func (m *Map[T, V]) Range(fn func(key T, value V) bool) {
	m.internalMap.Range(func(key, value any) bool {
		k := key.(T)
		v := value.(V)
		return fn(k, v)
	})
}

// Load ...
func (m *Map[T, V]) Load(key T) (value V, ok bool) {
	v, ok := m.internalMap.Load(key)
	if ok {
		value = v.(V)
	}

	return
}

// Delete ...
func (m *Map[T, V]) Delete(key T) {
	m.internalMap.Delete(key)
}

// NewMape ...
func NewMap[T any, V any]() *Map[T, V] {
	m := new(Map[T, V])
	m.internalMap = new(sync.Map)
	return m
}
