package safegroup

import (
	"sync"

	"github.com/WhiCu/async/group"
	"github.com/WhiCu/async/try"
)

type Group struct {
	wg sync.WaitGroup

	sem chan struct{}

	errOnce sync.Once
	err     error
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

func (g *Group) rawGo(f func()) {
	g.wg.Go(
		func() {
			if err := try.Try(f); err != nil {
				g.errOnce.Do(func() {
					g.err = err
				})
			}
			g.decrement()
		},
	)
}

func (g *Group) rawGoErr(f func() error) {
	g.wg.Go(
		func() {
			if err := try.TryErr(f); err != nil {
				g.errOnce.Do(func() {
					g.err = err
				})
			}
			g.decrement()
		},
	)
}

func (g *Group) Go(f func()) {
	g.increment()

	g.rawGo(f)

}

func (g *Group) GoErr(f func() error) {
	g.increment()

	g.rawGoErr(f)
}

func (g *Group) TryGo(f func()) error {
	if g.sem != nil {
		select {
		case g.sem <- struct{}{}:
		default:
			return group.ErrLimitExceeded
		}
	}

	g.rawGo(f)
	return nil
}

func (g *Group) TryGoErr(f func() error) error {
	if g.sem != nil {
		select {
		case g.sem <- struct{}{}:
		default:
			return group.ErrLimitExceeded
		}
	}

	g.rawGoErr(f)
	return nil
}

func (g *Group) Wait() error {
	g.wg.Wait()
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
