package ggkit

import (
	"context"

	"github.com/dsx137/gg-kit/internal/channel"
	"github.com/dsx137/gg-kit/internal/concurrent"
	"github.com/dsx137/gg-kit/internal/lang"
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

func NewAtomic[T any]() *Atomic[T]                    { return concurrent.NewAtomic[T]() }
func NewPool[T any](new func() T) *concurrent.Pool[T] { return concurrent.NewPool(new) }
func NewReentrantLock() *ReentrantLock                { return concurrent.NewReentrantLock() }
func NewReentrantRWLock() *ReentrantRWLock            { return concurrent.NewReentrantRWLock() }

// lang

func GetGoroutineId() int { return lang.GetGoroutineId() }
func Useless(v any)       { lang.Useless(v) }
