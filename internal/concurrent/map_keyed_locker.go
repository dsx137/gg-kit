package concurrent

import (
	"sync"

	"github.com/dsx137/gg-kit/internal/generic"
)

type MapKeyedLocker[K comparable] struct {
	locks *generic.SyncMap[K, *sync.RWMutex]
}

func NewMapKeyedLocker[K comparable]() *MapKeyedLocker[K] {
	return &MapKeyedLocker[K]{
		locks: generic.NewSyncMap[K, *sync.RWMutex](),
	}
}

func (k *MapKeyedLocker[K]) Lock(key K) func() {
	mu, _ := k.locks.LoadOrStore(key, &sync.RWMutex{})
	mu.Lock()

	return mu.Unlock
}

func (k *MapKeyedLocker[K]) TryLock(key K) (func(), bool) {
	mu, _ := k.locks.LoadOrStore(key, &sync.RWMutex{})
	if !mu.TryLock() {
		return nil, false
	}

	return mu.Unlock, true
}

func (k *MapKeyedLocker[K]) RLock(key K) func() {
	mu, _ := k.locks.LoadOrStore(key, &sync.RWMutex{})
	mu.RLock()

	return mu.RUnlock
}

func (k *MapKeyedLocker[K]) TryRLock(key K) (func(), bool) {
	mu, _ := k.locks.LoadOrStore(key, &sync.RWMutex{})
	if !mu.TryRLock() {
		return nil, false
	}

	return mu.RUnlock, true
}

func (k *MapKeyedLocker[K]) Locker(key K) sync.Locker {
	mu, _ := k.locks.LoadOrStore(key, &sync.RWMutex{})
	return mu
}

func (k *MapKeyedLocker[K]) RLocker(key K) sync.Locker {
	mu, _ := k.locks.LoadOrStore(key, &sync.RWMutex{})
	return mu.RLocker()
}
