package generic

import (
	"sync"
)

type SyncMap[K comparable, V any] struct {
	m *sync.Map
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{m: &sync.Map{}}
}

func (m *SyncMap[K, V]) Clear() {
	m.m.Clear()
}

func (m *SyncMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	return m.m.CompareAndSwap(key, old, new)
}

func (m *SyncMap[K, V]) CompareAndDelete(key K, old V) bool {
	return m.m.CompareAndDelete(key, old)
}

func (m *SyncMap[K, V]) Load(key K) (V, bool) {
	val, ok := m.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return val.(V), true
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	actual, loaded := m.m.LoadOrStore(key, value)
	return actual.(V), loaded
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.m.Delete(key)
}

func (m *SyncMap[K, V]) LoadAndDelete(key K) (V, bool) {
	actual, loaded := m.m.LoadAndDelete(key)
	if !loaded {
		var zero V
		return zero, false
	}
	return actual.(V), true
}

func (m *SyncMap[K, V]) Swap(key K, value V) (V, bool) {
	actual, swapped := m.m.Swap(key, value)
	if !swapped {
		var zero V
		return zero, false
	}
	return actual.(V), true
}

func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

// --------------- EXPAND ----------------

func (m *SyncMap[K, V]) All() func(yield func(K, V) bool) {
	return func(yield func(K, V) bool) {
		m.Range(func(k K, v V) bool {
			return yield(k, v)
		})
	}
}
