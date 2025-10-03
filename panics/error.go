package panics

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	// SkipFrames is the number of frames to skip when creating a panic error.
	SkipFrames = 3
)

type PanicError struct {
	Value   any
	Callers []uintptr
	Stack   []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic: %v", e.Value)
}

func (e *PanicError) Unwrap() error {
	if err, ok := e.Value.(error); ok {
		return err
	}
	return nil
}

// NewPanicError creates a panic error from the given value.
// It includes the calling stack frames and the current goroutine stack trace.
// It is safe to call on multiple goroutines concurrently.
// It can be unwrapped with errors.Unwrap to get the original error.
func NewPanicError(value any) error {
	var callers [64]uintptr
	n := runtime.Callers(SkipFrames, callers[:])
	return &PanicError{
		Value:   value,
		Callers: callers[:n],
		Stack:   debug.Stack(),
	}
}

func AsPanicError(err error) bool {
	var pe *PanicError
	return errors.As(err, &pe)
}
