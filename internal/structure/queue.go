package structure

import "github.com/dsx137/gg-kit/internal/generic"

type Queue[T any] struct {
	l *generic.List[T]
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		l: generic.NewList[T](),
	}
}

func (q *Queue[T]) Enqueue(v T) {
	q.l.PushBack(v)
}

func (q *Queue[T]) Dequeue() (v T, ok bool) {
	front := q.l.Front()
	if front == nil {
		var zero T
		return zero, false
	}
	q.l.Remove(front)
	return front.Value(), true
}

func (q *Queue[T]) Len() int {
	return q.l.Len()
}
