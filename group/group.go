package group

import "context"

type Runner interface {
	Go(func())
	TryGo(func()) error
}

type ErrRunner interface {
	GoErr(func() error)
	TryGoErr(func() error) error
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
	CtxGo(context.Context, func(context.Context))
	CtxTryGo(context.Context, func(context.Context)) error
}

type CtxErrRunner interface {
	CtxGoErr(context.Context, func(context.Context) error)
	CtxTryGoErr(context.Context, func(context.Context) error) error
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
