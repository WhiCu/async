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
		g, groupCtx := ctxgroup.WithContext(context.Background())

		Convey("It should run multiple goroutines with CtxGo()", func() {
			var counter atomic.Int32
			g.CtxGo(context.Background(), func(ctx context.Context) { counter.Add(1) })
			g.CtxGo(context.Background(), func(ctx context.Context) { counter.Add(1) })

			err := g.Wait()

			So(err, ShouldBeNil)
			So(counter.Load(), ShouldEqual, 2)
			So(context.Cause(groupCtx), ShouldEqual, context.Canceled)
		})

		Convey("It should handle panic correctly", func() {
			g.CtxGo(context.Background(), func(ctx context.Context) { panic("boom") })

			err := g.Wait()

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "boom")
			So(context.Cause(groupCtx).Error(), ShouldContainSubstring, "boom")
		})

		Convey("It should propagate error from CtxGoErr()", func() {
			testErr := errors.New("something went wrong")

			g.CtxGoErr(context.Background(), func(ctx context.Context) error {
				return testErr
			})

			err := g.Wait()

			So(err, ShouldEqual, testErr)
			So(context.Cause(groupCtx), ShouldEqual, testErr)
		})

		Convey("It should respect SetLimit()", func() {
			So(g.SetLimit(2), ShouldBeNil)

			called := atomic.Int32{}
			for i := 0; i < 2; i++ {
				err := g.CtxTryGo(context.Background(), func(ctx context.Context) {
					time.Sleep(10 * time.Millisecond)
					called.Add(1)
				})
				So(err, ShouldBeNil)
			}

			err := g.CtxTryGo(context.Background(), func(ctx context.Context) {})
			So(err, ShouldEqual, group.ErrLimitExceeded)

			So(g.Wait(), ShouldBeNil)
			So(called.Load(), ShouldEqual, 2)
		})

		Convey("It should return ErrModifyLimit if SetLimit called again after work started", func() {
			So(g.SetLimit(1), ShouldBeNil)
			g.CtxGo(context.Background(), func(ctx context.Context) { <-ctx.Done() })
			So(g.SetLimit(2), ShouldEqual, group.ErrModifyLimit)

			g.Cancel()
			So(g.Wait(), ShouldBeNil)
		})

		Convey("It should return ErrNegativeLimit for negative limit", func() {
			err := g.SetLimit(-1)
			So(err, ShouldEqual, group.ErrNegativeLimit)
		})
	})
}

func TestGroup_CancelPropagation(t *testing.T) {
	Convey("Given a group with context cancellation", t, func() {
		g, groupCtx := ctxgroup.WithContext(context.Background())

		var i atomic.Int32

		g.CtxGo(context.Background(), func(ctx context.Context) {
			select {
			case <-time.After(5 * time.Second):
				i.Store(1)
			case <-ctx.Done():
				i.Store(2)
			}
		})

		g.CtxGoErr(context.Background(), func(ctx context.Context) error {
			return errors.New("fatal")
		})

		err := g.Wait()

		So(err.Error(), ShouldEqual, "fatal")
		So(context.Cause(groupCtx).Error(), ShouldEqual, "fatal")
		So(i.Load(), ShouldEqual, int32(2))
	})
}

func TestGroup_IsolatedCancellation(t *testing.T) {
	Convey("Given two independent groups sharing the same parent context", t, func() {
		parent := context.Background()

		group1, ctx1 := ctxgroup.WithContext(parent)
		group2, ctx2 := ctxgroup.WithContext(parent)

		var g1Cancelled, g2Cancelled atomic.Bool

		group1.CtxGo(context.Background(), func(ctx context.Context) {
			<-ctx.Done()
			g1Cancelled.Store(true)
		})

		group2.CtxGo(context.Background(), func(ctx context.Context) {
			select {
			case <-ctx.Done():
				g2Cancelled.Store(true)
			case <-time.After(50 * time.Millisecond):
			}
		})

		group1.Cancel()
		So(group1.Wait(), ShouldBeNil)

		So(g1Cancelled.Load(), ShouldBeTrue)
		So(context.Cause(ctx1), ShouldEqual, context.Canceled)

		So(context.Cause(ctx2), ShouldBeNil)

		So(group2.Wait(), ShouldBeNil)
	})
}
