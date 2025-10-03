package panics

func Try(f func()) (err error) {
	var t Trier[any]
	return t.Try(f)
}

func TryValue[T any](f func() T) (v T, err error) {
	var t Trier[T]
	return t.TryValue(f)
}

func TryErr(f func() error) (err error) {
	var t Trier[any]
	return t.TryErr(f)
}

func TryValueErr[T any](f func() (T, error)) (v T, err error) {
	var t Trier[T]
	return t.TryValueErr(f)
}
