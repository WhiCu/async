package async_test

import (
	"fmt"
	"testing"

	"github.com/WhiCu/async"
)

func BenchmarkSafeGroup_NoPanic(b *testing.B) {
	var wg async.SafeGroup
	for _, workers := range []int{1, 4, 16} {
		b.Run(fmt.Sprintf("Workers_%d", workers), func(b *testing.B) {
			for i := 0; i < b.N; i++ {

				wg.SetLimit(workers)

				// b.ResetTimer()
				for j := 0; j < workers; j++ {
					wg.Go(func() {
						x := 0
						x++
					})
				}

				_ = wg.Wait()
			}
		})
	}
}

func BenchmarkSafeGroup_Panic(b *testing.B) {
	var wg async.SafeGroup
	for _, workers := range []int{1, 4, 16} {
		b.Run(fmt.Sprintf("Workers_%d", workers), func(b *testing.B) {
			for i := 0; i < b.N; i++ {

				wg.SetLimit(workers)

				// b.ResetTimer()
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
