// Package wsgatherer - this file contains logic for joining contexts
package wsgatherer

import (
	"context"
	"sync"
	"time"
)

type joinContext struct {
	mutex sync.Mutex
	ctx1  context.Context
	ctx2  context.Context
	done  chan struct{}
	err   error
}

//Join joins two contexts
func Join(ctx1, ctx2 context.Context) (context.Context, context.CancelFunc) {
	joined := &joinContext{ctx1: ctx1, ctx2: ctx2, done: make(chan struct{})}
	go joined.run()

	return joined, joined.cancel
}

func (joined *joinContext) Done() <-chan struct{} {
	return joined.done
}

func (joined *joinContext) Err() error {
	joined.mutex.Lock()
	defer joined.mutex.Unlock()

	return joined.err
}

func (joined *joinContext) Deadline() (deadline time.Time, ok bool) {
	d1, ok1 := joined.ctx1.Deadline()
	if !ok1 {
		return joined.ctx2.Deadline()
	}

	d2, ok2 := joined.ctx2.Deadline()
	if !ok2 {
		return d1, true
	}

	if d2.Before(d1) {
		return d2, true
	}

	return d1, true
}

func (joined *joinContext) Value(key interface{}) interface{} {
	v := joined.ctx1.Value(key)
	if v == nil {
		v = joined.ctx2.Value(key)
	}

	return v
}

//waits until any of contexts are done, sets error from it, closes joined done channel
func (joined *joinContext) run() {
	var doneCtx context.Context
	select {
	case <-joined.ctx1.Done():
		doneCtx = joined.ctx1
	case <-joined.ctx2.Done():
		doneCtx = joined.ctx2
	case <-joined.done:
		return
	}

	joined.mutex.Lock()
	defer joined.mutex.Unlock()

	if joined.err != nil {
		return
	}

	joined.err = doneCtx.Err()
	close(joined.done)
}

// sets error to context.Canceled, closes done channel
func (joined *joinContext) cancel() {
	joined.mutex.Lock()
	defer joined.mutex.Unlock()

	if joined.err != nil {
		return
	}

	joined.err = context.Canceled
	close(joined.done)
}
