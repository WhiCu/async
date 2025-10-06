package ctxgroup

import (
	"context"
	"sync"

	"github.com/WhiCu/async/try"
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

func (g *Group) Go(f func(context.Context)) {
	g.increment()

	g.wg.Go(
		func() {
			if err := try.Try(func() { f(g.Ctx) }); err != nil {
				g.errOnce.Do(func() {
					g.err = err
					g.Cancel()
				})
			}
			g.decrement()
		},
	)
}

func (g *Group) GoErr(f func(context.Context) error) {
	g.increment()

	g.wg.Go(
		func() {
			if err := try.TryErr(func() error { return f(g.Ctx) }); err != nil {
				g.errOnce.Do(func() {
					g.err = err
					g.Cancel()
				})
			}
			g.decrement()
		},
	)
}

func (g *Group) TryGo(f func(context.Context)) error {
	if g.sem != nil {
		select {
		case g.sem <- struct{}{}:
		default:
			return ErrLimitExceeded
		}
	}

	g.wg.Go(
		func() {
			if err := try.Try(func() { f(g.Ctx) }); err != nil {
				g.errOnce.Do(func() {
					g.err = err
					g.Cancel()
				})
			}
			g.decrement()
		},
	)
	return nil
}

func (g *Group) TryGoErr(f func(context.Context) error) error {
	if g.sem != nil {
		select {
		case g.sem <- struct{}{}:
		default:
			return ErrLimitExceeded
		}
	}

	g.wg.Go(
		func() {
			if err := try.TryErr(func() error { return f(g.Ctx) }); err != nil {
				g.errOnce.Do(func() {
					g.err = err
					g.Cancel()
				})
			}
			g.decrement()
		},
	)
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
		return ErrModifyLimit
	}

	if n < 0 {
		g.sem = nil
		return ErrNegativeLimit
	}
	g.sem = make(chan struct{}, n)
	return nil
}
