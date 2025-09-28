// Package panics provides utilities for safely handling panics in Go programs.
// It offers a Trier type that can execute functions and capture any panics that occur,
// converting them to errors for easier handling in Go's error-based control flow.
package panics

import (
	"sync/atomic"
)

type (
	// TrierAny is a type alias for Trier[any] for convenience.
	TrierAny = Trier[any]
)

// Trier provides thread-safe panic handling for functions.
// It can execute functions and capture any panics that occur,
// converting them to errors
type Trier[T any] struct {
	value atomic.Pointer[any]
	// worked atomic.Bool
}

// storePanic stores the panic value atomically with a stable heap-allocated address.
func (t *Trier[T]) storePanic(r any) {
	// copy value to a new variable so the address is stable/heap-allocated
	v := r
	t.value.Store(&v)
}

// recover captures any panic that occurs and stores it atomically.
func (t *Trier[T]) recover() {
	if r := recover(); r != nil {
		t.storePanic(r)
	}
}

// Try executes function f and returns an error if it panicked.
// If f panics, the panic is converted to an error and returned.
func (t *Trier[T]) Try(f func()) (err error) {
	t.Clean()
	func() {
		defer t.recover()
		f()
	}()
	return t.PanicAsError()
}

// TryValue executes function f and returns its result along with any panic as an error.
// If f panics, the zero value of T is returned and the panic is converted to an error.
func (t *Trier[T]) TryValue(f func() T) (v T, err error) {
	t.Clean()
	func() {
		defer t.recover()
		v = f()
	}()
	return v, t.PanicAsError()

}

// TryErr executes function f and returns its error, or a panic error if f panicked.
// If f panics, the panic is converted to an error and returned instead of the original error.
func (t *Trier[T]) TryErr(f func() error) (err error) {
	t.Clean()
	func() {
		defer t.recover()
		err = f()
	}()
	return t.сhangeErrorIfPanic(err)
}

// TryValueErr executes function f and returns its value, error, and any panic as an error.
// If f panics, zero values are returned for value and error, and the panic is converted to an error.
func (t *Trier[T]) TryValueErr(f func() (T, error)) (v T, err error) {
	t.Clean()
	func() {
		defer t.recover()
		v, err = f()
	}()
	return v, t.сhangeErrorIfPanic(err)
}

// Repanic re-throws the stored panic value, if any.
// If no panic was stored, it panics with nil.
func (t *Trier[T]) Repanic() {
	panic(t.Value())
}

// Value returns the stored panic value, or nil if no panic occurred.
func (t *Trier[T]) Value() any {
	if t.value.Load() == nil {
		return nil
	}
	return *t.value.Load()
}

// PanicAsError returns the stored panic value as an error, or nil if no panic occurred.
func (t *Trier[T]) PanicAsError() error {
	return AsPanicError(t.Value())
}

// сhangeErrorIfPanic returns the original error if no panic occurred, or converts the panic to an error.
func (t *Trier[T]) сhangeErrorIfPanic(err error) error {
	if t.Value() == nil {
		return err
	}
	return AsPanicError(t.Value())
}

// Clean resets the Trier state, clearing any stored panic information.
func (t *Trier[T]) Clean() {
	t.value.Store(nil)
}

// Worked returns true if the last executed function panicked.
// It uses CompareAndSwap to check if a panic value is stored.
func (t *Trier[T]) Worked() bool {
	return !t.value.CompareAndSwap(nil, nil)
}
