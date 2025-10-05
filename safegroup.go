package async

import (
	"github.com/WhiCu/async/try"
	"golang.org/x/sync/errgroup"
)

type SafeGroup struct {
	errgroup errgroup.Group
}

func (wg *SafeGroup) Go(f func()) {
	wg.errgroup.Go(func() error {
		return try.Try(f)
	})
}

func (wg *SafeGroup) Wait() error {
	return wg.errgroup.Wait()
}

func (wg *SafeGroup) SetLimit(n int) {
	wg.errgroup.SetLimit(n)
}
