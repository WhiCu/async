package safegroup

import (
	"github.com/WhiCu/async/try"
	"golang.org/x/sync/errgroup"
)

type Group struct {
	errgroup.Group
}

func (wg *Group) Go(f func()) {
	wg.Group.Go(func() error {
		return try.Try(f)
	})
}

func (wg *Group) GoErr(f func() error) {
	wg.Group.Go(func() error {
		return try.TryErr(f)
	})
}

func (wg *Group) TryGo(f func()) bool {
	return wg.Group.TryGo(func() error {
		return try.Try(f)
	})
}

func (wg *Group) TryGoErr(f func() error) bool {
	return wg.Group.TryGo(func() error {
		return try.TryErr(f)
	})
}

func (wg *Group) SetLimit(limit int) error {
	return try.Try(func() { wg.Group.SetLimit(limit) })
}
