package ggkit

import (
	"context"

	"github.com/dsx137/gg-kit/internal/channel"
	"github.com/dsx137/gg-kit/internal/concurrent"
	"github.com/dsx137/gg-kit/internal/lang"
	"github.com/dsx137/gg-kit/internal/structure"
	"github.com/dsx137/gg-kit/internal/util"

)

// channel
func Consume[T any](ch <-chan T, handler func(T) bool) { channel.Consume(ch, handler) }
func ConsumeWithCtx[T any](ctx context.Context, ch <-chan T, handler func(T) bool) {
	channel.ConsumeWithCtx(ctx, ch, handler)
}

// concurrent
type Atomic[T any] = concurrent.Atomic[T]
type Pool[T any] = concurrent.Pool[T]
type ReentrantLock = concurrent.ReentrantLock
type ReentrantRWLock = concurrent.ReentrantRWLock
type ReusePool[T any] = concurrent.ReusePool[T]

func NewAtomic[T any]() *Atomic[T]         { return concurrent.NewAtomic[T]() }
func NewPool[T any](new func() T) *Pool[T] { return concurrent.NewPool(new) }
func NewReentrantLock() *ReentrantLock     { return concurrent.NewReentrantLock() }
func NewReentrantRWLock() *ReentrantRWLock { return concurrent.NewReentrantRWLock() }
func NewReusePoolWithParentCtx[T any](parentCtx context.Context, factory func() (*T, error), validator func(*T) bool, closer func(*T) error) (*ReusePool[T], error) {
	return concurrent.NewReusePoolWithParentCtx(parentCtx, factory, validator, closer)
}
func NewReusePool[T any](factory func() (*T, error), validator func(*T) bool, closer func(*T) error) (*ReusePool[T], error) {
	return concurrent.NewReusePool(factory, validator, closer)
}

// lang
func GetGoroutineId() int { return lang.GetGoroutineId() }
func Useless(v any)       { lang.Useless(v) }

// structure
type List[T any] = structure.List[T]
type Queue[T any] = structure.Queue[T]

func NewList[T any]() *List[T]   { return structure.NewList[T]() }
func NewQueue[T any]() *Queue[T] { return structure.NewQueue[T]() }

// util
func GenerateRandomBytes(n int) ([]byte, error)    { return util.GenerateRandomBytes(n) }
func GenerateBase64Key(length int) (string, error) { return util.GenerateBase64Key(length) }
func GenerateHexKey(length int) (string, error)    { return util.GenerateHexKey(length) }
