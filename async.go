// Package async provides utilities for safe goroutine execution with panic handling.
// It wraps the panics package functionality for easier async operations.
package async

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/WhiCu/async/panics"
)

// Goroutine represents a running goroutine with panic-safe execution.
type Goroutine struct {
	done   chan struct{}
	err    error
	panic  any
	ctx    context.Context
	cancel context.CancelFunc
	once   sync.Once
}

// SafeGo runs a function in a goroutine with panic handling.
// Returns a Goroutine handle that allows waiting for completion and error retrieval.
func SafeGo(f func()) *Goroutine {
	g := &Goroutine{
		done: make(chan struct{}),
	}
	g.ctx, g.cancel = context.WithCancel(context.Background())

	go func() {
		defer g.once.Do(g.finish)

		var t panics.TrierAny
		err := t.Try(f)
		if err != nil {
			g.err = err
			g.panic = t.Value()
		}
	}()

	return g
}

// SafeGoWithContext runs a function in a goroutine with context support.
// The goroutine can be cancelled through the context.
func SafeGoWithContext(ctx context.Context, f func()) *Goroutine {
	g := &Goroutine{
		done: make(chan struct{}),
	}
	g.ctx, g.cancel = context.WithCancel(ctx)

	go func() {
		defer g.once.Do(g.finish)

		select {
		case <-ctx.Done():
			g.err = ctx.Err()
			return
		default:
			var t panics.TrierAny
			err := t.Try(f)
			if err != nil {
				g.err = err
				g.panic = t.Value()
			}
		}
	}()

	return g
}

// SafeGoWithTimeout runs a function with a timeout.
// If the timeout is reached, the operation is cancelled.
func SafeGoWithTimeout(timeout time.Duration, f func()) *Goroutine {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return SafeGoWithContext(ctx, f)
}

// Go runs a function in a goroutine with panic handling, returning a channel for results.
func Go[T any](f func() T) <-chan Result[T] {
	result := make(chan Result[T], 1)

	go func() {
		defer close(result)

		var t panics.Trier[T]
		value, err := t.TryValue(f)
		result <- Result[T]{
			Value: value,
			Error: err,
			Panic: t.Value(),
		}
	}()

	return result
}

// GoErr runs a function that returns a value and error in a goroutine.
func GoErr[T any](f func() (T, error)) <-chan Result[T] {
	result := make(chan Result[T], 1)

	go func() {
		defer close(result)

		var t panics.Trier[T]
		value, err := t.TryValueErr(f)
		result <- Result[T]{
			Value: value,
			Error: err,
			Panic: t.Value(),
		}
	}()

	return result
}

// GoWithCallback runs a function with a callback for handling results.
func GoWithCallback[T any](f func() T, callback func(T, error)) {
	go func() {
		var t panics.Trier[T]
		value, err := t.TryValue(f)
		callback(value, err)
	}()
}

// Result represents the result of a goroutine execution.
type Result[T any] struct {
	Value T     // The return value
	Error error // Any error or panic error
	Panic any   // The original panic value if panic occurred
}

// finish closes the done channel to signal completion.
func (g *Goroutine) finish() {
	close(g.done)
}

// Wait blocks until the goroutine completes and returns any error.
func (g *Goroutine) Wait() error {
	<-g.done
	return g.err
}

// WaitFor blocks for a maximum duration and returns an error if timeout occurs.
func (g *Goroutine) WaitFor(timeout time.Duration) error {
	select {
	case <-g.done:
		return g.err
	case <-time.After(timeout):
		return fmt.Errorf("goroutine timeout after %v", timeout)
	}
}

// Panic returns the original panic value if one occurred.
func (g *Goroutine) Panic() any {
	return g.panic
}

// HasPanicked returns true if the goroutine panicked.
func (g *Goroutine) HasPanicked() bool {
	return g.panic != nil
}

// Cancel cancels the goroutine execution.
func (g *Goroutine) Cancel() {
	g.cancel()
}

// Done returns a channel that signals when the goroutine is finished.
func (g *Goroutine) Done() <-chan struct{} {
	return g.done
}

// Parallel runs multiple functions in parallel goroutines and waits for all to complete.
func Parallel(functions ...func()) []error {
	var wg sync.WaitGroup
	results := make([]error, len(functions))

	for i, fn := range functions {
		wg.Add(1)
		go func(index int, f func()) {
			defer wg.Done()
			var t panics.TrierAny
			results[index] = t.Try(f)
		}(i, fn)
	}

	wg.Wait()
	return results
}

// ParallelWithResults runs multiple functions in parallel and returns all results.
func ParallelWithResults[T any](functions ...func() T) []Result[T] {
	var wg sync.WaitGroup
	results := make([]Result[T], len(functions))

	for i, fn := range functions {
		wg.Add(1)
		go func(index int, f func() T) {
			defer wg.Done()
			var t panics.Trier[T]
			value, err := t.TryValue(f)
			results[index] = Result[T]{
				Value: value,
				Error: err,
				Panic: t.Value(),
			}
		}(i, fn)
	}

	wg.Wait()
	return results
}

// RetryWithBackoff runs a function with exponential backoff retry on failures.
func RetryWithBackoff[T any](attempts int, backoff time.Duration, f func() (T, error)) (T, error) {
	var lastErr error
	var result T

	for i := 0; i < attempts; i++ {
		var t panics.Trier[T]
		result, lastErr = t.TryValueErr(f)

		if lastErr == nil {
			return result, nil
		}

		if i < attempts-1 {
			time.Sleep(backoff * time.Duration(1<<i))
		}
	}

	return result, fmt.Errorf("failed after %d attempts, last error: %w", attempts, lastErr)
}

// Pool manages a pool of goroutines for concurrent execution with limits.
type Pool struct {
	semaphore chan struct{}
	wg        sync.WaitGroup
}

// NewPool creates a new goroutine pool with the specified size limit.
func NewPool(size int) *Pool {
	return &Pool{
		semaphore: make(chan struct{}, size),
	}
}

// Submit submits a task to the pool for execution.
func (p *Pool) Submit(f func()) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		// Acquire semaphore
		p.semaphore <- struct{}{}
		defer func() { <-p.semaphore }()

		// Execute with panic handling
		var t panics.TrierAny
		t.Try(f)
	}()
}

// Wait waits for all submitted tasks to complete.
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Close closes the pool and waits for remaining tasks.
func (p *Pool) Close() {
	p.Wait()
	close(p.semaphore)
}
