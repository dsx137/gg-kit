package concurrent

import (
	"sync"
)

type ShardedKeyedLocker[K comparable] struct {
	shards []sync.RWMutex
	mask   uint64
	hash   func(K) uint64
}

func NewShardedKeyedLocker[K comparable](exp uint, hash func(K) uint64) *ShardedKeyedLocker[K] {
	if exp < 1 || exp > 32 {
		panic("exp must be between 1 and 32")
	}
	shardCount := 1 << exp
	return &ShardedKeyedLocker[K]{
		shards: make([]sync.RWMutex, shardCount),
		mask:   uint64(shardCount - 1),
		hash:   hash,
	}
}

func (k *ShardedKeyedLocker[K]) Lock(key K) func() {
	mu := &k.shards[k.hash(key)&k.mask]
	mu.Lock()
	return mu.Unlock
}

func (k *ShardedKeyedLocker[K]) TryLock(key K) (func(), bool) {
	mu := &k.shards[k.hash(key)&k.mask]
	if !mu.TryLock() {
		return nil, false
	}
	return mu.Unlock, true
}

func (k *ShardedKeyedLocker[K]) RLock(key K) func() {
	mu := &k.shards[k.hash(key)&k.mask]
	mu.RLock()
	return mu.RUnlock
}

func (k *ShardedKeyedLocker[K]) TryRLock(key K) (func(), bool) {
	mu := &k.shards[k.hash(key)&k.mask]
	if !mu.TryRLock() {
		return nil, false
	}
	return mu.RUnlock, true
}

func (k *ShardedKeyedLocker[K]) Locker(key K) sync.Locker {
	mu := &k.shards[k.hash(key)&k.mask]
	return mu
}

func (k *ShardedKeyedLocker[K]) RLocker(key K) sync.Locker {
	mu := &k.shards[k.hash(key)&k.mask]
	return mu.RLocker()
}
