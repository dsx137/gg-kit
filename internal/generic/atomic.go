package generic

import "sync/atomic"

type Atomic[T any] struct {
	a *atomic.Value
}

func NewAtomic[T any]() *Atomic[T] {
	return &Atomic[T]{a: &atomic.Value{}}
}

func NewAtomicWithValue[T any](val T) *Atomic[T] {
	ret := NewAtomic[T]()
	ret.Store(val)
	return ret
}

func (a *Atomic[T]) Store(val T)                    { a.a.Store(val) }
func (a *Atomic[T]) Load() T                        { return a.a.Load().(T) }
func (a *Atomic[T]) Swap(val T) T                   { return a.a.Swap(val).(T) }
func (a *Atomic[T]) CompareAndSwap(old, new T) bool { return a.a.CompareAndSwap(old, new) }
