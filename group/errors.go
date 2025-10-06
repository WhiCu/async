package group

import (
	"errors"
)

var (
	ErrLimitExceeded = errors.New("group: limit exceeded")
	ErrModifyLimit   = errors.New("group: modify limit while goroutines in the group are still active")
	ErrNegativeLimit = errors.New("group: negative limit")
)
