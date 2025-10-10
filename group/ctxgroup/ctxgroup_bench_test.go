package ctxgroup_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/WhiCu/async/group/ctxgroup"
)

func BenchmarkCtxGroup_NoCtx(b *testing.B) {
	data := []byte("The quick brown fox jumps over the lazy dog")
	wg, _ := ctxgroup.WithContext(context.Background())
	numGo := 1024
	for b.Loop() {
		for i := 0; i < numGo; i++ {
			wg.CtxGo(context.Background(), func(ctx context.Context) {
				for i := 0; i < 1024; i++ {
					sha256.Sum256(data)
				}

			})
		}
		_ = wg.Wait()
	}
}
func BenchmarkCtxGroup_Ctx(b *testing.B) {
	data := []byte("The quick brown fox jumps over the lazy dog")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg, _ := ctxgroup.WithContext(ctx)
	numGo := 1024
	for b.Loop() {
		for i := 0; i < numGo; i++ {
			wg.CtxGo(ctx, func(ctx context.Context) {
				for i := 0; i < 1024; i++ {
					sha256.Sum256(data)
				}

			})
		}
		_ = wg.Wait()
	}
}

func BenchmarkCtxGroup_NoPanic_NoCtx(b *testing.B) {

	for _, workers := range []int{1, 4, 16} {
		wg, _ := ctxgroup.WithContext(context.Background())
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

func BenchmarkCtxGroup_NoPanic(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, workers := range []int{1, 4, 16} {
		wg, _ := ctxgroup.WithContext(context.Background())
		_ = wg.SetLimit(workers)

		b.Run(fmt.Sprintf("Workers_%d", workers), func(b *testing.B) {
			var number atomic.Int64
			for i := 0; i < b.N; i++ {
				for j := 0; j < workers; j++ {
					wg.CtxGo(ctx, func(ctx context.Context) {
						number.Add(1)
					})
				}
				_ = wg.Wait()
			}
		})
	}
}

func BenchmarkCtxGroup_Panic(b *testing.B) {
	for _, workers := range []int{1, 4, 16} {
		wg, _ := ctxgroup.WithContext(context.Background())
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
