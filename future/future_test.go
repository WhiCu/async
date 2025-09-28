package future

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFuture(t *testing.T) {
	Convey("Given a Promise without error handling", t, func() {
		Convey("When the function returns a value successfully", func() {
			fut := Promise(func() int { return 42 })
			val, err := fut.Value()

			Convey("Then it should return the correct value and no panic", func() {
				So(val, ShouldEqual, 42)
				So(err, ShouldBeNil)
				So(fut.ValuePanic(), ShouldBeNil)
			})
		})

		Convey("When the function panics", func() {
			fut := Promise(func() int {
				panic("boom")
			})
			val, err := fut.Value()

			Convey("Then it should detect the panic and still return zero value", func() {
				So(val, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				So(fut.ValuePanic(), ShouldEqual, "boom")
			})
		})
	})

	Convey("Given a PromiseErr with error handling", t, func() {
		Convey("When the function returns a value without error", func() {
			fut := PromiseErr(func() (int, error) {
				return 7, nil
			})
			val, err := fut.Value()

			Convey("Then it should return the value and no error or panic", func() {
				So(val, ShouldEqual, 7)
				So(err, ShouldBeNil)
				So(fut.ValuePanic(), ShouldBeNil)
			})
		})

		Convey("When the function returns an error", func() {
			fut := PromiseErr(func() (int, error) {
				return 0, errors.New("fail")
			})
			val, err := fut.Value()

			Convey("Then it should return the error and no panic", func() {
				So(val, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "fail")
				So(fut.ValuePanic(), ShouldBeNil)
			})
		})

		Convey("When the function panics", func() {
			fut := PromiseErr(func() (int, error) {
				panic("err-panic")
			})
			val, err := fut.Value()

			Convey("Then it should detect the panic and still return zero value and nil error", func() {
				So(val, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "err-panic")
				So(fut.ValuePanic(), ShouldEqual, "err-panic")
			})
		})
	})
}
