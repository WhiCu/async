package try

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTry(t *testing.T) {
	Convey("Given the Try functions", t, func() {

		Convey("Try with a safe function", func() {
			err := Try(func() {})
			Convey("Then it should not return an error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("Try with a panicking function", func() {
			err := Try(func() { panic("boom") })
			Convey("Then it should capture the panic as an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "boom")
			})
		})

		Convey("TryValue with a safe function", func() {
			v, err := TryValue(func() int { return 42 })
			Convey("Then it should return the value and no error", func() {
				So(v, ShouldEqual, 42)
				So(err, ShouldBeNil)
			})
		})

		Convey("TryValue with a panicking function", func() {
			_, err := TryValue(func() int { panic("fail") })
			Convey("Then it should capture the panic as an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "fail")
			})
		})

		Convey("TryErr with a safe function returning nil", func() {
			err := TryErr(func() error { return nil })
			Convey("Then it should return nil error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("TryErr with a function returning an error", func() {
			err := TryErr(func() error { return errors.New("fail") })
			Convey("Then it should return that error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "fail")
			})
		})

		Convey("TryErr with a panicking function", func() {
			err := TryErr(func() error { panic("panic in err") })
			Convey("Then it should capture the panic as an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "panic in err")
			})
		})

		Convey("TryValueErr with a safe function", func() {
			v, err := TryValueErr(func() (int, error) { return 99, nil })
			Convey("Then it should return value and no error", func() {
				So(v, ShouldEqual, 99)
				So(err, ShouldBeNil)
			})
		})

		Convey("TryValueErr with a function returning an error", func() {
			v, err := TryValueErr(func() (int, error) { return 0, errors.New("bad") })
			Convey("Then it should return zero value and the error", func() {
				So(v, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "bad")
			})
		})

		Convey("TryValueErr with a panicking function", func() {
			_, err := TryValueErr(func() (int, error) { panic("panic in valueErr") })
			Convey("Then it should capture the panic as an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "panic in valueErr")
			})
		})
	})
}
