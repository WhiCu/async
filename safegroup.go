package async

import (
	"github.com/WhiCu/async/try"
	"golang.org/x/sync/errgroup"
)

type SafeGroup struct {
	errgroup.Group
}

func (wg *SafeGroup) Go(f func()) {
	wg.Group.Go(func() error {
		return try.Try(f)
	})
}

func (wg *SafeGroup) GoErr(f func() error) {
	wg.Group.Go(func() error {
		return try.TryErr(f)
	})
}

func (wg *SafeGroup) TryGo(f func()) bool {
	return wg.Group.TryGo(func() error {
		return try.Try(f)
	})
}

func (wg *SafeGroup) TryGoErr(f func() error) bool {
	return wg.Group.TryGo(func() error {
		return try.TryErr(f)
	})
}
