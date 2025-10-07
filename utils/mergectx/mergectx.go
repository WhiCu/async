package mergectx

import (
	"context"
	"time"
)

type mergeContext struct {
	primary, secondary context.Context
	ch                 chan struct{}
}

func MergeContext(primary context.Context, secondary context.Context) context.Context {
	tc := mergeContext{
		primary:   primary,
		secondary: secondary,
		ch:        make(chan struct{}),
	}
	go mergeChan(tc.ch, primary.Done(), secondary.Done())
	return &tc
}

func mergeChan(out chan<- struct{}, firstIn, secondIn <-chan struct{}) {
	select {
	case <-firstIn:
	case <-secondIn:
	}
	close(out)
}

func (c *mergeContext) Done() <-chan struct{} {
	return c.ch
}

func (c *mergeContext) Err() error {
	if err := c.primary.Err(); err != nil {
		return err
	}
	return c.secondary.Err()
}

func (c *mergeContext) Deadline() (deadline time.Time, ok bool) {
	d1, ok1 := c.primary.Deadline()
	d2, ok2 := c.secondary.Deadline()

	if !ok1 {
		return d2, ok2
	}
	if !ok2 {
		return d1, ok1
	}

	if d1.Before(d2) {
		return d1, true
	}
	return d2, true
}

func (c *mergeContext) Value(key interface{}) interface{} {
	if v := c.primary.Value(key); v != nil {
		return v
	}
	return c.secondary.Value(key)
}
