package concurrent

import "sync"

func WithLock(locker sync.Locker, f func()) {
	locker.Lock()
	defer locker.Unlock()
	f()
}

func WithLockResult[T any](locker sync.Locker, f func() T) T {
	locker.Lock()
	defer locker.Unlock()
	return f()
}

func WithLockResultAndError[T any](locker sync.Locker, f func() (T, error)) (T, error) {
	locker.Lock()
	defer locker.Unlock()
	return f()
}
