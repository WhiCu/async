package async

import (
	"github.com/WhiCu/async/try"
)

func GoTry[T any](f func()) {
	go func() {
		try.Try(f)
	}()
}
