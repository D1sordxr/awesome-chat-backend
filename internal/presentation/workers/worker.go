package workers

import (
	"awesome-chat/internal/domain/app/ports"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

type Handler interface {
	Start(ctx context.Context) error
}

type Worker struct {
	log      ports.Logger
	handlers []Handler
	errChan  chan error
	done     chan struct{}
}

func NewWorker(
	log ports.Logger,
	handlers ...Handler,
) *Worker {
	return &Worker{
		log:      log,
		handlers: handlers,
		errChan:  make(chan error),
		done:     make(chan struct{}),
	}
}

func (w *Worker) Start(ctx context.Context) error {
	w.log.Info("starting worker", "total_handlers", len(w.handlers))

	errGroup, gCtx := errgroup.WithContext(ctx)
	go func() {
		w.errChan <- errGroup.Wait()
		close(w.done)
	}()

	for idx, handler := range w.handlers {
		func(idx int, handler Handler) {
			errGroup.Go(func() error {
				w.log.Info("starting worker handler", "idx", idx)
				return handler.Start(gCtx)
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
	select {
	case <-w.done:
		w.log.Info("all handlers stopped")
		return nil
	case <-ctx.Done():
		w.log.Warn("forced shutdown due to context timeout")
		return ctx.Err()
	}
}
