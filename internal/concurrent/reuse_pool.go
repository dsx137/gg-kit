package concurrent

import (
	"sync"

	"github.com/dsx137/gg-kit/internal/structure"
)

type ReusePool[T any] struct {
	mu        *sync.Mutex
	resources *structure.Queue[*T]
	factory   func() (*T, error)
	validator func(*T) bool
	closer    func(*T) error
}

func NewReusePool[T any](factory func() (*T, error), validator func(*T) bool, closer func(*T) error) (*ReusePool[T], error) {
	pool := &ReusePool[T]{
		mu:        &sync.Mutex{},
		resources: structure.NewQueue[*T](),
		factory:   factory,
		validator: validator,
		closer:    closer,
	}
	return pool, nil
}

func (p *ReusePool[T]) Get() (*T, error) {
	for {
		p.mu.Lock()
		e, ok := p.resources.Dequeue()
		p.mu.Unlock()
		if !ok {
			if p.factory != nil {
				return p.factory()
			}
			return nil, nil
		}
		if p.validator == nil || p.validator(e) {
			return e, nil
		}
		if p.closer == nil {
			continue
		}
		_ = p.closer(e)
	}
}

func (p *ReusePool[T]) Put(res *T) error {
	if res == nil {
		return nil
	}

	if p.validator != nil && !p.validator(res) {
		if p.closer != nil {
			_ = p.closer(res)
		}
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.resources.Enqueue(res)
	return nil
}

func (p *ReusePool[T]) Clear() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closer == nil {
		p.resources = structure.NewQueue[*T]()
		return nil
	}

	var firstErr error
	for {
		res, ok := p.resources.Dequeue()
		if !ok {
			break
		}
		if err := p.closer(res); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
