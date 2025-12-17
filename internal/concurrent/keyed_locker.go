package concurrent

import "sync"

type KeyedLocker[K comparable] interface {
	Lock(key K) func()
	TryLock(key K) func()
	RLock(key K) func()
	TryRLock(key K) func()
	Locker(key K) sync.Locker
	RLocker(key K) sync.Locker
}
