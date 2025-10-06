package ctxgroup

import (
	"errors"
)

var (
	ErrLimitExceeded = errors.New("ctxgroup: limit exceeded")
	ErrCanceled      = errors.New("ctxgroup: canceled")
	ErrModifyLimit   = errors.New("ctxgroup: modify limit while goroutines in the group are still active")
	ErrNegativeLimit = errors.New("ctxgroup: negative limit")
)
