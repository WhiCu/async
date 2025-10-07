package mergectx

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type ctxKey string

const (
	keyFoo ctxKey = "foo"
	keyBar ctxKey = "bar"
)

func TestTeeContext(t *testing.T) {
	Convey("Given two contexts", t, func() {

		Convey("When neither is cancelled", func() {
			ctx1, cancel1 := context.WithCancel(context.Background())
			ctx2, cancel2 := context.WithCancel(context.Background())
			defer cancel1()
			defer cancel2()

			tc := MergeContext(ctx1, ctx2)

			Convey("Then Done channel should not be closed yet", func() {
				select {
				case <-tc.Done():
					So("should not be closed", ShouldEqual, "")
				default:
					So(true, ShouldBeTrue)
				}
			})
		})

		Convey("When the primary context is cancelled", func() {
			ctx1, cancel1 := context.WithCancel(context.Background())
			ctx2 := context.WithoutCancel(context.Background())

			tc := MergeContext(ctx1, ctx2)

			cancel1()

			Convey("Then Done channel should close", func() {
				select {
				case <-tc.Done():
					So(true, ShouldBeTrue)
				case <-time.After(50 * time.Millisecond):
					So("timeout", ShouldEqual, "")
				}
			})

			Convey("And Err should match primary", func() {
				<-tc.Done()
				So(tc.Err(), ShouldEqual, context.Canceled)
			})
		})

		Convey("When the secondary context is cancelled", func() {
			ctx1 := context.WithoutCancel(context.Background())
			ctx2, cancel2 := context.WithCancel(context.Background())

			tc := MergeContext(ctx1, ctx2)

			cancel2()

			Convey("Then Done channel should close", func() {
				select {
				case <-tc.Done():
					So(true, ShouldBeTrue)
				case <-time.After(50 * time.Millisecond):
					So("timeout", ShouldEqual, "")
				}
			})

			Convey("And Err should match secondary", func() {
				<-tc.Done()
				So(tc.Err(), ShouldEqual, context.Canceled)
			})
		})

		Convey("When both have deadlines", func() {
			ctx1, cancel1 := context.WithTimeout(context.Background(), 100*time.Millisecond)
			ctx2, cancel2 := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel1()
			defer cancel2()

			tc := MergeContext(ctx1, ctx2)

			Convey("Then TeeContext should return the earliest deadline", func() {
				d, ok := tc.Deadline()
				So(ok, ShouldBeTrue)
				So(time.Until(d), ShouldBeLessThanOrEqualTo, 150*time.Millisecond)
			})
		})

		Convey("When one has no deadline", func() {
			ctx1 := context.Background()
			ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
			defer cancel2()

			tc := MergeContext(ctx1, ctx2)

			Convey("Then TeeContext should use the secondary deadline", func() {
				_, ok := tc.Deadline()
				So(ok, ShouldBeTrue)
			})
		})

		Convey("When using Value lookup", func() {
			ctx1 := context.WithValue(context.Background(), keyFoo, "foo")
			ctx2 := context.WithValue(context.Background(), keyBar, "bar")

			tc := MergeContext(ctx1, ctx2)

			Convey("Then it should find in primary first", func() {
				So(tc.Value(keyFoo), ShouldEqual, "foo")
			})

			Convey("And fallback to secondary", func() {
				So(tc.Value(keyBar), ShouldEqual, "bar")
			})

			Convey("And return nil for missing key", func() {
				So(tc.Value("nope"), ShouldBeNil)
			})
		})
	})
}

func TestMergeContext_Cause(t *testing.T) {
	Convey("Given two contexts with cancel cause", t, func() {
		ctx1, cancel1 := context.WithCancelCause(context.Background())
		ctx2, cancel2 := context.WithCancelCause(context.Background())
		defer cancel1(nil)
		defer cancel2(nil)

		tc := MergeContext(ctx1, ctx2)

		Convey("When neither is cancelled", func() {
			So(context.Cause(tc), ShouldBeNil)
		})

		Convey("When primary is cancelled with a cause", func() {
			cause := context.DeadlineExceeded
			cancel1(cause)

			select {
			case <-tc.Done():
			case <-time.After(50 * time.Millisecond):
				So("timeout", ShouldEqual, "")
			}

			Convey("Then Cause() should return primary cause", func() {
				So(context.Cause(tc), ShouldEqual, cause)
			})
		})

		Convey("When secondary is cancelled with a cause", func() {
			cause := context.Canceled
			cancel2(cause)

			select {
			case <-tc.Done():
			case <-time.After(50 * time.Millisecond):
				So("timeout", ShouldEqual, "")
			}

			Convey("Then Cause() should return secondary cause", func() {
				So(context.Cause(tc), ShouldEqual, cause)
			})
		})
	})
}
