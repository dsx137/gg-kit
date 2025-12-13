package generic

import "container/list"

type Element[T any] struct {
	le *list.Element
}

func (e *Element[T]) Value() T {
	return e.le.Value.(T)
}

func (e *Element[T]) Next() *Element[T] {
	return &Element[T]{le: e.le.Next()}
}
func (e *Element[T]) Prev() *Element[T] {
	return &Element[T]{le: e.le.Prev()}
}

type List[T any] struct {
	l *list.List
}

func NewList[T any]() *List[T] {
	return &List[T]{
		l: list.New(),
	}
}

func (l *List[T]) Back() *Element[T] {
	le := l.l.Back()
	if le == nil {
		return nil
	}
	return &Element[T]{le: le}
}

func (l *List[T]) Front() *Element[T] {
	le := l.l.Front()
	if le == nil {
		return nil
	}
	return &Element[T]{le: le}
}

func (l *List[T]) Init() *List[T] {
	l.l.Init()
	return l
}

func (l *List[T]) InsertAfter(v T, mark *Element[T]) *Element[T] {
	return &Element[T]{le: l.l.InsertAfter(v, mark.le)}
}

func (l *List[T]) InsertBefore(v T, mark *Element[T]) *Element[T] {
	return &Element[T]{le: l.l.InsertBefore(v, mark.le)}
}

func (l *List[T]) Len() int {
	return l.l.Len()
}

func (l *List[T]) MoveAfter(e *Element[T], mark *Element[T]) {
	l.l.MoveAfter(e.le, mark.le)
}

func (l *List[T]) MoveBefore(e *Element[T], mark *Element[T]) {
	l.l.MoveBefore(e.le, mark.le)
}

func (l *List[T]) MoveToBack(e *Element[T]) {
	l.l.MoveToBack(e.le)
}

func (l *List[T]) MoveToFront(e *Element[T]) {
	l.l.MoveToFront(e.le)
}

func (l *List[T]) PushBack(v T) *Element[T] {
	return &Element[T]{le: l.l.PushBack(v)}
}

func (l *List[T]) PushBackList(other *List[T]) {
	l.l.PushBackList(other.l)
}

func (l *List[T]) PushFront(v T) *Element[T] {
	return &Element[T]{le: l.l.PushFront(v)}
}

func (l *List[T]) PushFrontList(other *List[T]) {
	l.l.PushFrontList(other.l)
}

func (l *List[T]) Remove(e *Element[T]) T {
	if e == nil {
		var zero T
		return zero
	}
	return l.l.Remove(e.le).(T)
}

// --------------- EXPAND ----------------

func (l *List[T]) All() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		for e := l.Front(); e != nil; e = e.Next() {
			if !yield(e.Value()) {
				break
			}
		}
	}
}

func (l *List[T]) ReverseAll() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		for e := l.Back(); e != nil; e = e.Prev() {
			if !yield(e.Value()) {
				break
			}
		}
	}
}
