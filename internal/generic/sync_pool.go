package generic

import "sync"

type SyncPool[T any] struct {
	p *sync.Pool
}

func NewSyncPool[T any](new func() T) *SyncPool[T] {
	return &SyncPool[T]{
		p: &sync.Pool{New: func() any { return new() }},
	}
}

func (p *SyncPool[T]) Get() T  { return p.p.Get().(T) }
func (p *SyncPool[T]) Put(x T) { p.p.Put(x) }
