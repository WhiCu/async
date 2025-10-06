package group

import "context"

type Group interface {
	Go(f func())
	GoErr(f func() error)
	TryGo(f func()) bool
	TryGoErr(f func() error) bool
	SetLimit(limit int) error
	Wait() error
}

type CtxGroup interface {
	Group
	CtxGo(f func(context.Context))
	CtxGoErr(f func(context.Context) error)
	CtxTryGo(f func(context.Context)) bool
	CtxTryGoErr(f func(context.Context) error) bool
	Cancel()
}
