package panics

import "fmt"

type PanicError struct {
	Panic any
}

func (e PanicError) Error() string {
	return fmt.Sprintf("Panic: %v", e.Panic)
}

func AsPanicError(value any) error {
	if value == nil {
		return nil
	}
	return PanicError{Panic: value}
}
