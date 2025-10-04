package panics

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTrier(t *testing.T) {
	Convey("Given a new Trier", t, func() {
		tr := &Trier[int]{}

		Convey("When running a function that does not panic", func() {
			err := tr.Try(func() {})

			Convey("Then it should not detect a panic", func() {
				So(err, ShouldBeNil)
				So(tr.Worked(), ShouldBeFalse)
				So(tr.Value(), ShouldBeNil)
			})
		})

		Convey("When running a function that panics", func() {
			err := tr.Try(func() {
				panic("boom")
			})

			Convey("Then it should detect the panic", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "boom")
				So(tr.Worked(), ShouldBeTrue)
				So(tr.Value(), ShouldEqual, "boom")
			})

			Convey("And when Repanic is called", func() {
				Convey("Then it should rethrow the stored panic", func() {
					So(func() { tr.Repanic() }, ShouldPanicWith, "boom")
				})
			})

			Convey("And when Clean is called", func() {
				tr.Clean()

				Convey("Then the stored panic value should be nil", func() {
					So(tr.Value(), ShouldBeNil)
					So(tr.Worked(), ShouldBeFalse)
				})
			})
		})

		Convey("When using TryValue with a safe function", func() {
			v, err := tr.TryValue(func() int { return 42 })

			Convey("Then it should return the value and no panic", func() {
				So(v, ShouldEqual, 42)
				So(err, ShouldBeNil)
			})
		})

		Convey("When using TryValue with a panicking function", func() {
			_, err := tr.TryValue(func() int { panic("fail") })

			Convey("Then it should detect the panic", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "fail")
			})
		})

		Convey("When using TryErr with a safe function", func() {
			err := tr.TryErr(func() error { return nil })

			Convey("Then it should return no error and no panic", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When using TryErr with a failing function", func() {
			err := tr.TryErr(func() error { return errors.New("fail") })

			Convey("Then it should return the error and no panic", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "fail")
				So(tr.Worked(), ShouldBeFalse)
			})
		})

		Convey("When using TryErr with a panicking function", func() {
			err := tr.TryErr(func() error { panic("panic in err") })

			Convey("Then it should detect the panic", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "panic in err")
				So(tr.Worked(), ShouldBeTrue)
			})
		})

		Convey("When using TryValueErr with a safe function", func() {
			v, err := tr.TryValueErr(func() (int, error) {
				return 99, nil
			})

			Convey("Then it should return the value, no error, and no panic", func() {
				So(v, ShouldEqual, 99)
				So(err, ShouldBeNil)
				So(tr.Worked(), ShouldBeFalse)
			})
		})

		Convey("When using TryValueErr with an error function", func() {
			v, err := tr.TryValueErr(func() (int, error) {
				return 0, errors.New("bad")
			})

			Convey("Then it should return zero value, the error, and no panic", func() {
				So(v, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "bad")
				So(tr.Worked(), ShouldBeFalse)
			})
		})

		Convey("When using TryValueErr with a panicking function", func() {
			_, err := tr.TryValueErr(func() (int, error) {
				panic("panic in valueErr")
			})

			Convey("Then it should detect the panic", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "panic in valueErr")
				So(tr.Worked(), ShouldBeTrue)
			})
		})
	})
}
