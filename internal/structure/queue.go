package structure

type Queue[T any] struct {
	l *List[T]
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		l: NewList[T](),
	}
}

func (q *Queue[T]) Enqueue(v T) {
	q.l.PushBack(v)
}

func (q *Queue[T]) Dequeue() (v T, ok bool) {
	front := q.l.Front()
	if front.e == nil {
		var zero T
		return zero, false
	}
	q.l.Remove(front)
	return front.Value(), true
}

func (q *Queue[T]) Len() int {
	return q.l.Len()
}
