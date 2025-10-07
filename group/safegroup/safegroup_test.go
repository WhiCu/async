package safegroup_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/WhiCu/async/group"
	"github.com/WhiCu/async/group/safegroup"
	"github.com/WhiCu/async/try"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSafeGroup(t *testing.T) {
	Convey("Given a SafeGroup instance", t, func() {
		var sg safegroup.Group

		Convey("It should run multiple successful goroutines", func() {
			var count atomic.Int32
			for i := 0; i < 5; i++ {
				sg.Go(func() { count.Add(1) })
			}

			err := sg.Wait()
			So(err, ShouldBeNil)
			So(count.Load(), ShouldEqual, 5)
		})

		Convey("It should capture panic from a goroutine", func() {
			sg.Go(func() { panic("boom") })
			sg.Go(func() {})

			err := sg.Wait()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "boom")
			So(try.AsPanicError(err), ShouldBeTrue)
		})

		Convey("It should respect SetLimit()", func() {
			So(sg.SetLimit(1), ShouldBeNil)

			var counter atomic.Int32
			start := time.Now()

			for i := 0; i < 3; i++ {
				sg.Go(func() {
					time.Sleep(50 * time.Millisecond)
					counter.Add(1)
				})
			}

			err := sg.Wait()
			So(err, ShouldBeNil)
			So(counter.Load(), ShouldEqual, 3)

			elapsed := time.Since(start)
			So(elapsed, ShouldBeGreaterThanOrEqualTo, 150*time.Millisecond)
		})

		Convey("It should return ErrLimitExceeded from TryGo()", func() {
			So(sg.SetLimit(1), ShouldBeNil)

			done := make(chan struct{})
			err := sg.TryGo(func() {
				time.Sleep(100 * time.Millisecond)
				close(done)
			})
			So(err, ShouldBeNil)

			err = sg.TryGo(func() {})
			So(err, ShouldEqual, group.ErrLimitExceeded)

			<-done
			So(sg.Wait(), ShouldBeNil)
		})

		Convey("It should propagate explicit errors from TryGoErr()", func() {
			testErr := errors.New("custom error")
			err := sg.TryGoErr(func() error { return testErr })
			So(err, ShouldBeNil)

			waitErr := sg.Wait()
			So(waitErr, ShouldEqual, testErr)
		})

		Convey("It should return ErrModifyLimit if SetLimit is called after work started", func() {
			So(sg.SetLimit(1), ShouldBeNil)
			sg.Go(func() { time.Sleep(50 * time.Millisecond) })

			err := sg.SetLimit(2)
			So(err, ShouldEqual, group.ErrModifyLimit)

			So(sg.Wait(), ShouldBeNil)
		})

		Convey("It should return ErrNegativeLimit for negative limit", func() {
			err := sg.SetLimit(-1)
			So(err, ShouldEqual, group.ErrNegativeLimit)
		})
	})
}
