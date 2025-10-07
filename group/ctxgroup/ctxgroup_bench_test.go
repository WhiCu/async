package ctxgroup_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/WhiCu/async/group/ctxgroup"
)

func BenchmarkCtxGroup_NoPanic(b *testing.B) {
	var wg ctxgroup.Group
	for _, workers := range []int{1, 4, 16} {
		_ = wg.SetLimit(workers)
		b.Run(fmt.Sprintf("Workers_%d", workers), func(b *testing.B) {
			var number atomic.Int64
			for i := 0; i < b.N; i++ {
				for j := 0; j < workers; j++ {
					wg.CtxGo(context.Background(), func(ctx context.Context) {
						number.Add(1)
					})
				}
				_ = wg.Wait()
			}
		})
	}
}

func BenchmarkCtxGroup_Panic(b *testing.B) {
	var wg ctxgroup.Group
	for _, workers := range []int{1, 4, 16} {
		_ = wg.SetLimit(workers)
		b.Run(fmt.Sprintf("Workers_%d", workers), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < workers; j++ {
					wg.CtxGo(context.Background(), func(ctx context.Context) {
						panic("boom")
					})
				}
				_ = wg.Wait()
			}
		})
	}
}

func BenchmarkCtxGroup_ContextCancel(b *testing.B) {
	for _, workers := range []int{1, 4, 16} {
		b.Run(fmt.Sprintf("Workers_%d", workers), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ctx, cancel := context.WithCancel(context.Background())
				wg, _ := ctxgroup.WithContext(ctx)
				_ = wg.SetLimit(workers)

				var cancelled atomic.Int64
				for j := 0; j < workers; j++ {
					wg.CtxGo(context.Background(), func(ctx context.Context) {
						<-ctx.Done()
						cancelled.Add(1)
					})
				}

				cancel()
				_ = wg.Wait()

				if got := cancelled.Load(); got != int64(workers) {
					b.Fatalf("expected %d cancelled goroutines, got %d", workers, got)
				}
			}
		})
	}
}
