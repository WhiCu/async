package stream

import "iter"

type Iterator[T any] struct {
	iter iter.Seq[T]
}

func (i Iterator[T]) Slice() []T {
	collect := make([]T, 0)
	for v := range i.iter {
		collect = append(collect, v)
	}
	return collect
}

func From[T any](slice []T) *Iterator[T] {
	return &Iterator[T]{
		iter: func(yield func(T) bool) {
			for _, v := range slice {
				if !yield(v) {
					break
				}
			}
		},
	}
}

func (i Iterator[T]) Each(f func(T)) *Iterator[T] {
	return &Iterator[T]{
		iter: func(yield func(T) bool) {
			for v := range i.iter {
				f(v)
				if !yield(v) {
					break
				}
			}
		},
	}

}

func (i Iterator[T]) GoEach(f func(T)) Iterator[T] {
	for v := range i {
		f(v)
	}
	return i
}

// utils

func copy(i Iterator[T], iter iter.Seq[T]) *Iterator[T] {
	i.iter = iter
	return i
}
