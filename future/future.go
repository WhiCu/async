// Package future provides asynchronous execution utilities with error-based panic handling.
// It allows running functions in goroutines and converts any panics to errors for easier handling.
package future

import "github.com/WhiCu/async/panics"

// future represents an asynchronous computation that converts panics to errors.
// It wraps a channel for receiving results and uses Trier for panic handling.
type future[T any] struct {
	value <-chan valueError[T]

	t panics.Trier[T]
}

// valueError wraps a value and error result from a future computation.
type valueError[T any] struct {
	value T
	err   error
}

// Promise creates a future that executes function f asynchronously.
// It returns a future that will contain the result or an error if the function panicked.
func Promise[T any](f func() T) *future[T] {
	c := make(chan valueError[T])
	future := &future[T]{
		value: c,
	}
	go func() {
		defer close(c)
		v, err := future.t.TryValue(f)
		c <- valueError[T]{
			value: v,
			err:   err,
		}
	}()
	return future
}

// PromiseErr creates a future that executes function f asynchronously and handles errors.
// It returns a future that will contain the result, error, or panic converted to error.
func PromiseErr[T any](f func() (T, error)) *future[T] {
	c := make(chan valueError[T])
	future := &future[T]{
		value: c,
	}
	go func() {
		defer close(c)
		v, err := future.t.TryValueErr(f)
		c <- valueError[T]{
			value: v,
			err:   err,
		}
	}()
	return future
}

// Value blocks until the future computation completes and returns the result.
// It returns the computed value and any error (including panics converted to errors).
func (f *future[T]) Value() (value T, err error) {
	r := <-f.value
	return r.value, r.err
}

// ValuePanic returns the original panic value that occurred during computation, if any.
// This provides access to the raw panic value before it was converted to an error.
func (f *future[T]) ValuePanic() any {
	return f.t.Value()
}
