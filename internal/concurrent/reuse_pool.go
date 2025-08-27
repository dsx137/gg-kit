package concurrent

import (
	"context"
	"sync"

	"github.com/dsx137/gg-kit/internal/structure"
)

type ReusePool[T any] struct {
	mu        *sync.Mutex
	resources *structure.Queue[*T]
	factory   func() (*T, error)
	validator func(*T) bool
	closer    func(*T) error
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewReusePoolWithParentCtx[T any](parentCtx context.Context, factory func() (*T, error), validator func(*T) bool, closer func(*T) error) (*ReusePool[T], error) {
	ctx, cancel := context.WithCancel(parentCtx)
	pool := &ReusePool[T]{
		mu:        &sync.Mutex{},
		resources: structure.NewQueue[*T](),
		factory:   factory,
		validator: validator,
		closer:    closer,
		ctx:       ctx,
		ctxCancel: cancel,
	}
	return pool, nil
}

func NewReusePool[T any](factory func() (*T, error), validator func(*T) bool, closer func(*T) error) (*ReusePool[T], error) {
	return NewReusePoolWithParentCtx(context.Background(), factory, validator, closer)
}

func (p *ReusePool[T]) Get() (*T, error) {
	if p.ctx.Err() != nil {
		return nil, p.ctx.Err()
	}

	var (
		res *T
		err error
	)

	for {
		p.mu.Lock()
		if p.resources.Len() <= 0 {
			p.mu.Unlock()
			break
		}
		e, ok := p.resources.Dequeue()
		p.mu.Unlock()
		if !ok {
			break
		}
		if p.validator(e) {
			res = e
			break
		}
		if err := p.closer(e); err != nil {
			return nil, err
		}
	}
	if res == nil {
		res, err = p.factory()
	}

	return res, err
}

func (p *ReusePool[T]) Put(res *T) error {
	if res == nil {
		return nil
	}

	if p.ctx.Err() != nil || !p.validator(res) {
		return p.closer(res)
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.resources.Enqueue(res)
	return nil
}

func (p *ReusePool[T]) Close() error {
	p.ctxCancel()

	p.mu.Lock()
	defer p.mu.Unlock()

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

func (p *ReusePool[T]) IsClosed() bool {
	return p.ctx.Err() != nil
}
