package async_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/WhiCu/async"
	"github.com/WhiCu/async/try"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSafeGroup(t *testing.T) {
	Convey("Given a SafeGroup instance", t, func() {
		var sg async.SafeGroup

		Convey("When running multiple successful goroutines", func() {
			var count atomic.Int32
			for i := 0; i < 5; i++ {
				sg.Go(func() { count.Add(1) })
			}

			err := sg.Wait()

			Convey("Then no error should occur", func() {
				So(err, ShouldBeNil)
			})

			Convey("And all goroutines should have completed", func() {
				So(count.Load(), ShouldEqual, 5)
			})
		})

		Convey("When one of the goroutines panics", func() {
			sg.Go(func() { panic("boom") })
			sg.Go(func() {}) // нормальная горутина

			err := sg.Wait()

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When SetLimit is applied", func() {
			sg.SetLimit(1)

			start := time.Now()

			for i := 0; i < 3; i++ {
				sg.Go(func() {
					time.Sleep(100 * time.Millisecond)
				})
			}

			sg.Wait()
			elapsed := time.Since(start)

			Convey("Then execution time should reflect the concurrency limit", func() {
				So(elapsed, ShouldBeGreaterThanOrEqualTo, 300*time.Millisecond)
			})
		})

		Convey("When a goroutine returns an explicit error", func() {
			expected := errors.New("custom error")
			sg.Go(func() {
				panic(expected)
			})

			err := sg.Wait()

			Convey("Then the returned error should match the panic value", func() {
				So(err.Error(), ShouldContainSubstring, "custom error")
			})
			Convey("And the error should be a PanicError", func() {
				So(try.AsPanicError(err), ShouldBeTrue)
			})
		})
	})
}
