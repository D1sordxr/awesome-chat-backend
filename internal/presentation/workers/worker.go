package workers

import (
	"awesome-chat/internal/domain/app/ports"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync/atomic"
	"time"
)

type Handler interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Worker struct {
	log      ports.Logger
	handlers []Handler
	
	errChan  chan error
	isClosed atomic.Bool
}

func NewWorker(
	log ports.Logger,
	handlers ...Handler,
) *Worker {
	return &Worker{
		log:      log,
		handlers: handlers,
		errChan:  make(chan error),
		isClosed: atomic.Bool{},
	}
}

func (w *Worker) Start(ctx context.Context) error {
	w.log.Info("starting worker", "total_handlers", len(w.handlers))

	errGroup, ctx := errgroup.WithContext(ctx)
	go func() {
		w.errChan <- errGroup.Wait()
	}()

	for idx, handler := range w.handlers {
		func(idx int, handler Handler) {
			errGroup.Go(func() error {
				return handler.Start(ctx)
			})
		}(idx, handler)
	}

	select {
	case err := <-w.errChan:
		w.log.Error("worker received critical error, initiating shutdown", "error", err)
		return fmt.Errorf("handler error: %w", err)
	case <-ctx.Done():
		return nil
	}
}

func (w *Worker) Shutdown(ctx context.Context) error {
	defer w.isClosed.Store(true)
	if w.isClosed.Load() {
		return nil
	}
	done := make(chan struct{})
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	go func() {
		defer close(done)
		for i := len(w.handlers) - 1; i >= 0; i-- {
			if err := w.handlers[i].Stop(ctx); err != nil {
				w.log.Error("handler stopped with error", "error", err.Error())
			}
		}
	}()

	select {
	case <-done:
		w.log.Info("all handlers stopped")
		return nil
	case <-ctx.Done():
		w.log.Warn("forced shutdown due to context timeout")
		return ctx.Err()
	}
}
