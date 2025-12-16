package concurrent

import (
	"sync"

	"github.com/dsx137/gg-kit/internal/generic"
)

type KeyedLocker[K comparable] struct {
	locks *generic.SyncMap[K, *sync.RWMutex]
}

func NewKeyedLocker[K comparable]() *KeyedLocker[K] {
	return &KeyedLocker[K]{
		locks: generic.NewSyncMap[K, *sync.RWMutex](),
	}
}

func (k *KeyedLocker[K]) Lock(key K) func() {
	mu, _ := k.locks.LoadOrStore(key, &sync.RWMutex{})
	mu.Lock()

	return mu.Unlock
}

func (k *KeyedLocker[K]) TryLock(key K) (func(), bool) {
	mu, _ := k.locks.LoadOrStore(key, &sync.RWMutex{})
	if !mu.TryLock() {
		return nil, false
	}

	return mu.Unlock, true
}

func (k *KeyedLocker[K]) RLock(key K) func() {
	mu, _ := k.locks.LoadOrStore(key, &sync.RWMutex{})
	mu.RLock()

	return mu.RUnlock
}

func (k *KeyedLocker[K]) TryRLock(key K) (func(), bool) {
	mu, _ := k.locks.LoadOrStore(key, &sync.RWMutex{})
	if !mu.TryRLock() {
		return nil, false
	}

	return mu.RUnlock, true
}
