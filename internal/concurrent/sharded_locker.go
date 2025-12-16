package concurrent

import (
	"sync"
)

type ShardedLocker[K comparable] struct {
	shards []sync.RWMutex
	mask   uint64
	hash   func(K) uint64
}

func NewShardedLocker[K comparable](exp uint, hash func(K) uint64) *ShardedLocker[K] {
	if exp < 1 || exp > 32 {
		panic("exp must be between 1 and 32")
	}
	shardCount := 1 << exp
	return &ShardedLocker[K]{
		shards: make([]sync.RWMutex, shardCount),
		mask:   uint64(shardCount - 1),
		hash:   hash,
	}
}

func (k *ShardedLocker[K]) Lock(key K) func() {
	mu := &k.shards[k.hash(key)&k.mask]
	mu.Lock()
	return mu.Unlock
}

func (k *ShardedLocker[K]) TryLock(key K) (func(), bool) {
	mu := &k.shards[k.hash(key)&k.mask]
	if !mu.TryLock() {
		return nil, false
	}
	return mu.Unlock, true
}

func (k *ShardedLocker[K]) RLock(key K) func() {
	mu := &k.shards[k.hash(key)&k.mask]
	mu.RLock()
	return mu.RUnlock
}

func (k *ShardedLocker[K]) TryRLock(key K) (func(), bool) {
	mu := &k.shards[k.hash(key)&k.mask]
	if !mu.TryRLock() {
		return nil, false
	}
	return mu.RUnlock, true
}
