package stream

import (
	"iter"

	"github.com/WhiCu/async"
)

type Iterator[T any] iter.Seq[T]

func (i Iterator[T]) Slice() []T {
	collect := make([]T, 0)
	for v := range i {
		collect = append(collect, v)
	}
	return collect
}

func From[T any](slice []T) Iterator[T] {
	return func(yield func(T) bool) {
		for _, v := range slice {
			if !yield(v) {
				break
			}
		}
	}
}

func (i Iterator[T]) Each(f func(T)) {
	for v := range i {
		f(v)
	}
}

func (i Iterator[T]) GoEach(f func(T)) error {
	sg := async.SafeGroup{}
	for v := range i {
		sg.Go(func() { f(v) })
	}
	return sg.Wait()
}

func (i Iterator[T]) GoEachLimit(f func(T), limit int) error {
	sg := async.SafeGroup{}
	sg.SetLimit(limit)
	for v := range i {
		sg.Go(func() { f(v) })
	}
	return sg.Wait()
}

func Map[T, R any](i Iterator[T], f func(T) R) Iterator[R] {
	return func(yield func(R) bool) {
		for v := range i {
			if !yield(f(v)) {
				break
			}
		}
	}
}

// utils

// func copy(i Iterator[T], iter iter.Seq[T]) *Iterator[T] {
// 	i.iter = iter
// 	return i
// }
