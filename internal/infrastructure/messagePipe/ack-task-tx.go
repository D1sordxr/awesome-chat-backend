package messagePipe

import (
	"context"
	"sync"
	"time"
)

type AckPipeTx struct {
	wg sync.WaitGroup
}

func NewAckPipeTx() *AckPipeTx {
	return &AckPipeTx{
		wg: sync.WaitGroup{},
	}
}

func (a *AckPipeTx) Add() {
	a.wg.Add(1)
}

func (a *AckPipeTx) Confirm() {
	a.wg.Done()
}

func (a *AckPipeTx) Wait() {
	a.wg.Wait()
}

func (a *AckPipeTx) Start(_ context.Context) error {
	a.wg.Add(1)
	return nil
}

func (a *AckPipeTx) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	a.wg.Done()

	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
