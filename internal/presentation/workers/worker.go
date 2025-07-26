package workers

import (
	"awesome-chat/internal/domain/app/ports"
	"context"
	"fmt"
	"sync"
)

type Handler interface {
	Start(ctx context.Context) error
}

type Worker struct {
	log      ports.Logger
	wg       sync.WaitGroup
	errChan  chan error
	handlers []Handler
}

func NewWorker(
	log ports.Logger,
	handlers ...Handler,
) *Worker {
	return &Worker{
		log:      log,
		wg:       sync.WaitGroup{},
		errChan:  make(chan error),
		handlers: handlers,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	w.log.Info("starting worker", "total_handlers", len(w.handlers))

	w.wg.Add(len(w.handlers))
	for i, handler := range w.handlers {
		go func(idx int, h Handler) {
			defer w.wg.Done()

			select {
			case <-ctx.Done():
				w.log.Debug("handler skipped due to shutdown", "handler_index", idx)
				return
			default:
				w.log.Debug("starting handler", "handler_index", idx)
				if err := h.Start(ctx); err != nil {
					w.log.Error("handler failed", "error", err, "handler_index", idx)
					select {
					case w.errChan <- err:
					case <-ctx.Done():
					}
				}
			}
		}(i, handler)
	}

	select {
	case err := <-w.errChan:
		w.log.Error("worker received critical error, initiating shutdown", "error", err)
		cancel()
		return fmt.Errorf("handler error: %w", err)

	case <-ctx.Done():
		w.log.Info("worker context cancelled")
	}

	return nil
}

func (w *Worker) Shutdown(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
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
