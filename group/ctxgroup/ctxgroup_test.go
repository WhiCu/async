package ctxgroup_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/WhiCu/async/group"
	"github.com/WhiCu/async/group/ctxgroup"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGroup_Behavior(t *testing.T) {
	Convey("Given a new ctxgroup.Group with context", t, func() {
		g, ctx := ctxgroup.WithContext(context.Background())

		Convey("When running multiple goroutines with Go()", func() {
			var counter atomic.Int32
			g.CtxGo(func(ctx context.Context) {
				counter.Add(1)
			})
			g.CtxGo(func(ctx context.Context) {
				counter.Add(1)
			})

			err := g.Wait()

			Convey("Then all goroutines should finish successfully", func() {
				So(err, ShouldBeNil)
				So(counter.Load(), ShouldEqual, 2)
				So(context.Cause(ctx), ShouldEqual, context.Canceled)
			})
		})

		Convey("When a goroutine panics", func() {
			g.CtxGo(func(ctx context.Context) {
				panic("boom")
			})

			err := g.Wait()

			Convey("Then Wait should return a panic error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "boom")
				So(context.Cause(ctx).Error(), ShouldContainSubstring, "boom")
			})
		})

		Convey("When a goroutine returns an error using GoErr()", func() {
			testErr := errors.New("something went wrong")

			g.CtxGoErr(func(ctx context.Context) error {
				return testErr
			})

			err := g.Wait()

			Convey("Then the error should propagate", func() {
				So(err, ShouldEqual, testErr)
				So(context.Cause(ctx), ShouldEqual, testErr)
			})
		})

		Convey("When SetLimit is used", func() {
			err := g.SetLimit(2)

			Convey("Then limit should be applied without error", func() {
				So(err, ShouldBeNil)
			})

			Convey("When running more goroutines than the limit with TryGo()", func() {
				called := atomic.Int32{}
				for i := 0; i < 2; i++ {
					So(g.CtxTryGo(func(ctx context.Context) {
						time.Sleep(10 * time.Millisecond)
						called.Add(1)
					}), ShouldBeNil)
				}

				err := g.CtxTryGo(func(ctx context.Context) {
					time.Sleep(10 * time.Millisecond)
				})

				Convey("Then the third call should fail with ErrLimitExceeded", func() {
					So(err, ShouldEqual, group.ErrLimitExceeded)
				})

				Convey("And the others should execute", func() {
					So(g.Wait(), ShouldBeNil)
					So(called.Load(), ShouldEqual, 2)
				})
			})
		})

		Convey("When SetLimit called", func() {
			So(g.SetLimit(1), ShouldBeNil)
			g.CtxGo(func(ctx context.Context) { <-ctx.Done() })
			Convey("Then ErrModifyLimit should be returned for the second call", func() {
				So(g.SetLimit(2), ShouldEqual, group.ErrModifyLimit)
				Convey("But Wait should still work", func() {
					g.Cancel()
					So(g.Wait(), ShouldBeNil)
				})

			})

		})

		Convey("When SetLimit called with negative value", func() {
			err := g.SetLimit(-1)
			Convey("Then ErrNegativeLimit should be returned", func() {
				So(err, ShouldEqual, group.ErrNegativeLimit)
			})
		})
	})
}

func TestGroup_CancelPropagation(t *testing.T) {
	Convey("Given a group with context cancellation", t, func() {
		g, ctx := ctxgroup.WithContext(context.Background())

		Convey("When one goroutine cancels the context", func() {
			var i atomic.Int32

			g.CtxGo(func(ctx context.Context) {
				select {
				case <-time.After(5 * time.Second):
					i.Store(1)
				case <-ctx.Done():
					i.Store(2)
				}
			})

			g.CtxGoErr(func(ctx context.Context) error {
				return errors.New("fatal")
			})

			err := g.Wait()

			Convey("Then the Wait should return the first error and cancel the context", func() {
				So(err.Error(), ShouldEqual, "fatal")

				So(context.Cause(ctx).Error(), ShouldEqual, "fatal")

				So(i.Load(), ShouldEqual, int32(2))
			})
		})
	})
}
