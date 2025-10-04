package panics_test

import (
	"errors"
	"testing"

	"github.com/WhiCu/async/panics"
)

func BenchmarkTry_NoPanic(b *testing.B) {
	f := func() {}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = panics.Try(f)
	}
}

func BenchmarkTry_Panic(b *testing.B) {
	f := func() { panic("boom") }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = panics.Try(f)
	}
}

func BenchmarkTryValue_NoPanic(b *testing.B) {

	f := func() int { return 42 }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = panics.TryValue(f)
	}
}

func BenchmarkTryValue_Panic(b *testing.B) {
	f := func() int { panic("boom") }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = panics.TryValue(f)
	}
}

func BenchmarkTryErr_NoPanic(b *testing.B) {
	f := func() error { return errors.New("fail") }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = panics.TryErr(f)
	}
}

func BenchmarkTryErr_Panic(b *testing.B) {
	f := func() error { panic("boom") }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = panics.TryErr(f)
	}
}

func BenchmarkTryValueErr_NoPanic(b *testing.B) {
	f := func() (int, error) { return 42, errors.New("fail") }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = panics.TryValueErr(f)
	}
}

func BenchmarkTryValueErr_Panic(b *testing.B) {
	f := func() (int, error) { panic("boom") }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = panics.TryValueErr(f)
	}
}
