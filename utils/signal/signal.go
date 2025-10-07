package signal

import (
	"context"
	"os"
	"os/signal"

	"github.com/WhiCu/async/utils/mergectx"
)

func MergeSignal(ctx context.Context, signals ...os.Signal) context.Context {
	ctxSig, _ := signal.NotifyContext(context.Background(), signals...)
	return mergectx.MergeContext(ctx, ctxSig)
}
