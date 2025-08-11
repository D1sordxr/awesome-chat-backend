package messagePipe

import (
	"context"
	"sync"
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
	<-ctx.Done()

	a.wg.Done()
	a.wg.Wait()
	
	return nil
}
