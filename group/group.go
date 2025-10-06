package group

import "context"

type Runner interface {
	Go(f func())
	TryGo(f func()) error
}

type ErrRunner interface {
	GoErr(f func() error)
	TryGoErr(f func() error) error
}

type Limiter interface {
	SetLimit(limit int) error
}

type Waiter interface {
	Wait() error
}

type Group interface {
	Runner
	ErrRunner
	Limiter
	Waiter
}

type CtxRunner interface {
	CtxGo(f func(context.Context))
	CtxTryGo(f func(context.Context)) error
}

type CtxErrRunner interface {
	CtxGoErr(f func(context.Context) error)
	CtxTryGoErr(f func(context.Context) error) error
}

type Canceler interface {
	Cancel()
}

type CtxGroup interface {
	CtxRunner
	CtxErrRunner
	Limiter
	Waiter
	Canceler
}
