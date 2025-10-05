package try

func Try(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = NewPanicError(r)
		}
	}()
	f()
	return err
}

// TryValue executes f and returns its result along with a panic as an error, if any.
func TryValue[T any](f func() T) (v T, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = NewPanicError(r)
		}
	}()
	v = f()
	return v, err
}

// TryErr executes f and returns its error, or a panic as an error, if any.
func TryErr(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = NewPanicError(r)
		}
	}()
	err = f()
	return err
}

// TryValueErr executes f and returns its value, error, and a panic as an error, if any.
func TryValueErr[T any](f func() (T, error)) (v T, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = NewPanicError(r)
		}
	}()
	v, err = f()
	return v, err
}
