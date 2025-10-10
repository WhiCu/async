package safegroup_test

import (
	"crypto/sha256"
	"fmt"
	"sync/atomic"
	"testing"

	safe "github.com/WhiCu/async/group/safegroup"
)

func BenchmarkSafeGroup(b *testing.B) {
	data := []byte("The quick brown fox jumps over the lazy dog")
	var wg safe.Group
	numGo := 1024
	for b.Loop() {
		for i := 0; i < numGo; i++ {
			wg.Go(func() {
				for i := 0; i < 1024; i++ {
					sha256.Sum256(data)
				}

			})
		}
		_ = wg.Wait()
	}
}

func BenchmarkSafeGroup_NoPanic(b *testing.B) {
	var wg safe.Group
	for _, workers := range []int{1, 4, 16} {
		wg.SetLimit(workers)
		b.Run(fmt.Sprintf("Workers_%d", workers), func(b *testing.B) {
			var number atomic.Int64
			for i := 0; i < b.N; i++ {
				for j := 0; j < workers; j++ {
					wg.Go(func() {
						number.Add(1)
					})
				}

				_ = wg.Wait()
			}
		})
	}
}

func BenchmarkSafeGroup_Panic(b *testing.B) {
	var wg safe.Group
	for _, workers := range []int{1, 4, 16} {
		wg.SetLimit(workers)
		b.Run(fmt.Sprintf("Workers_%d", workers), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < workers; j++ {
					wg.Go(func() {
						panic("boom")
					})
				}

				_ = wg.Wait()
			}
		})
	}
}
