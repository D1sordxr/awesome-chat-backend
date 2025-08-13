package batchSaver

import (
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/message/ports/store"
	"awesome-chat/internal/domain/core/message/ports/worker"
	"awesome-chat/internal/domain/core/message/vo"
	"context"
	"fmt"
	"time"
)

const (
	batchSize     = 64
	flushInterval = time.Second * 3
)

type Handler struct {
	log           appPorts.Logger
	ackWritePipe  worker.MessagePipe[string]
	saverReadPipe worker.MessagePipe[vo.StreamMessage]
	saverStore    store.SaveFromStreamStore
	ackPipeTx     worker.AckPipeTx
	batchSize     int
	flushInterval time.Duration
}

func NewHandler(
	log appPorts.Logger,
	ackPipe worker.MessagePipe[string],
	saverPipe worker.MessagePipe[vo.StreamMessage],
	saverStore store.SaveFromStreamStore,
	ackPipeTx worker.AckPipeTx,
) *Handler {
	return &Handler{
		log:           log,
		ackWritePipe:  ackPipe,
		saverReadPipe: saverPipe,
		saverStore:    saverStore,
		ackPipeTx:     ackPipeTx,
		batchSize:     batchSize,
		flushInterval: flushInterval,
	}
}

func (h *Handler) Start(ctx context.Context) error {
	const op = "message.batchSaver.Handler.Start"
	withFields := func(args ...any) []any {
		return append([]any{"operation", op}, args...)
	}

	h.log.Info("Starting batchSaver...", withFields()...)
	defer h.log.Info("BatchSaver stopped", withFields()...)

	ticker := time.NewTicker(h.flushInterval)
	defer ticker.Stop()

	batch := make([]vo.StreamMessage, 0, h.batchSize)

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}

		if err := h.saverStore.SaveBatch(ctx, batch); err != nil {
			return fmt.Errorf("%s: batch save failed: %w", op, err)
		}

		for _, msg := range batch {
			select {
			case h.ackWritePipe.GetWriteChan() <- msg.AckID:
				h.ackPipeTx.Add()
			case <-ctx.Done():
				return nil
			}
		}

		batch = batch[:0]

		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return flush()
		case <-ticker.C:
			if err := flush(); err != nil {
				h.log.Error("Failed to flush batch", withFields("error", err.Error())...)
			}
		case msg := <-h.saverReadPipe.GetReadChan():
			batch = append(batch, msg)

			if len(batch) >= h.batchSize {
				if err := flush(); err != nil {
					h.log.Error("Failed to flush batch", withFields("error", err.Error())...)
				}
			}
		}
	}
}

func (h *Handler) Stop(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		h.ackWritePipe.Close()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
