package ctxgroup

import (
	"context"
	"sync"

	"github.com/WhiCu/async/group"
	"github.com/WhiCu/async/try"
	"github.com/WhiCu/async/utils/mergectx"
)

type Group struct {
	cancel context.CancelCauseFunc
	Ctx    context.Context

	wg sync.WaitGroup

	sem chan struct{}

	errOnce sync.Once
	err     error
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	return &Group{Ctx: ctx, cancel: cancel}, ctx
}

func (g *Group) increment() {
	if g.sem != nil {
		g.sem <- struct{}{}
	}
}

func (g *Group) decrement() {
	if g.sem != nil {
		<-g.sem
	}
}

func mergeCtx(primary, secondary context.Context) context.Context {
	return mergectx.MergeContext(primary, secondary)
}

func (g *Group) rawGo(f func(context.Context), ctx context.Context) {
	g.wg.Go(
		func() {
			if err := try.Try(func() { f(ctx) }); err != nil {
				g.errOnce.Do(func() {
					g.err = err
					g.Cancel()
				})
			}
			g.decrement()
		},
	)
}

func (g *Group) rawGoErr(f func(context.Context) error, ctx context.Context) {
	g.wg.Go(
		func() {
			if err := try.TryErr(func() error { return f(ctx) }); err != nil {
				g.errOnce.Do(func() {
					g.err = err
					g.Cancel()
				})
			}
			g.decrement()
		},
	)
}

func (g *Group) CtxGo(ctx context.Context, f func(context.Context)) {
	g.increment()

	switch {
	case ctx == nil || ctx.Done() == nil:
		ctx = g.Ctx
	default:
		ctx = mergeCtx(g.Ctx, ctx)
	}
	g.rawGo(f, ctx)

}

func (g *Group) CtxGoErr(ctx context.Context, f func(context.Context) error) {
	g.increment()

	switch {
	case ctx == nil || ctx.Done() == nil:
		ctx = g.Ctx
	default:
		ctx = mergeCtx(g.Ctx, ctx)
	}
	g.rawGoErr(f, ctx)

}

func (g *Group) CtxTryGo(ctx context.Context, f func(context.Context)) error {
	if g.sem != nil {
		select {
		case g.sem <- struct{}{}:
		default:
			return group.ErrLimitExceeded
		}
	}

	switch {
	case ctx == nil || ctx.Done() == nil:
		ctx = g.Ctx
	default:
		ctx = mergeCtx(g.Ctx, ctx)
	}
	g.rawGo(f, ctx)

	return nil
}

func (g *Group) CtxTryGoErr(ctx context.Context, f func(context.Context) error) error {
	if g.sem != nil {
		select {
		case g.sem <- struct{}{}:
		default:
			return group.ErrLimitExceeded
		}
	}

	switch {
	case ctx == nil || ctx.Done() == nil:
		ctx = g.Ctx
	default:
		ctx = mergeCtx(g.Ctx, ctx)
	}
	g.rawGoErr(f, ctx)

	return nil
}

func (g *Group) Cancel() {
	if g.cancel != nil {
		switch err := g.err; err {
		case nil:
			g.cancel(context.Canceled)
		default:
			g.cancel(err)
		}
	}
}

func (g *Group) Wait() error {
	g.wg.Wait()
	g.Cancel()
	return g.err
}

func (g *Group) SetLimit(n int) error {
	if len(g.sem) != 0 {
		return group.ErrModifyLimit
	}

	if n < 0 {
		g.sem = nil
		return group.ErrNegativeLimit
	}

	g.sem = make(chan struct{}, n)
	return nil
}
