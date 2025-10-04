package panics_test

import (
	"errors"
	"testing"

	"github.com/WhiCu/async/panics"
)

func BenchmarkTrier_Try_NoPanic(b *testing.B) {
	t := &panics.Trier[int]{}
	f := func() {}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.Try(f)
	}
}

func BenchmarkTrier_Try_Panic(b *testing.B) {
	t := &panics.Trier[int]{}
	f := func() { panic("boom") }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.Try(f)
	}
}

func BenchmarkTrier_TryValue_NoPanic(b *testing.B) {
	t := &panics.Trier[int]{}
	f := func() int { return 42 }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = t.TryValue(f)
	}
}

func BenchmarkTrier_TryValue_Panic(b *testing.B) {
	t := &panics.Trier[int]{}
	f := func() int { panic("boom") }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = t.TryValue(f)
	}
}

func BenchmarkTrier_TryErr_NoPanic(b *testing.B) {
	t := &panics.Trier[int]{}
	f := func() error { return errors.New("fail") }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.TryErr(f)
	}
}

func BenchmarkTrier_TryErr_Panic(b *testing.B) {
	t := &panics.Trier[int]{}
	f := func() error { panic("boom") }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.TryErr(f)
	}
}

func BenchmarkTrier_TryValueErr_NoPanic(b *testing.B) {
	t := &panics.Trier[int]{}
	f := func() (int, error) { return 42, errors.New("fail") }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = t.TryValueErr(f)
	}
}

func BenchmarkTrier_TryValueErr_Panic(b *testing.B) {
	t := &panics.Trier[int]{}
	f := func() (int, error) { panic("boom") }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = t.TryValueErr(f)
	}
}
