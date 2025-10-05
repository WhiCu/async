package async_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/WhiCu/async"
)

func BenchmarkSafeGroup_NoPanic(b *testing.B) {
	var wg async.SafeGroup
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
	var wg async.SafeGroup
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
