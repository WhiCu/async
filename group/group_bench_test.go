package group_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"golang.org/x/sync/errgroup"
)

func BenchmarkErrGroup(b *testing.B) {
	var wg errgroup.Group
	for _, workers := range []int{1, 4, 16} {
		wg.SetLimit(workers)
		b.Run(fmt.Sprintf("Workers_%d", workers), func(b *testing.B) {
			var number atomic.Int64
			for b.Loop() {
				for j := 0; j < workers; j++ {
					wg.Go(func() error {
						number.Add(1)
						return nil
					})
				}

				_ = wg.Wait()
			}
		})
	}
}
