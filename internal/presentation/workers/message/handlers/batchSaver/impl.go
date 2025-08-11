package batchSaver

import (
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/message/ports"
	"awesome-chat/internal/domain/core/message/ports/worker"
	"awesome-chat/internal/domain/core/message/vo"
	"context"
	"fmt"
	"time"
)

const (
	batchSize     = 256
	flushInterval = time.Second * 3
)

type Handler struct {
	log           appPorts.Logger
	ackWritePipe  worker.MessagePipe[string]
	saverReadPipe worker.MessagePipe[vo.StreamMessage]
	repo          ports.Repository
	ackPipeTx     worker.AckPipeTx
	batchSize     int
	flushInterval time.Duration
}

func NewHandler(
	log appPorts.Logger,
	ackPipe worker.MessagePipe[string],
	saverPipe worker.MessagePipe[vo.StreamMessage],
	repo ports.Repository,
	ackPipeTx worker.AckPipeTx,
) *Handler {
	return &Handler{
		log:           log,
		ackWritePipe:  ackPipe,
		saverReadPipe: saverPipe,
		repo:          repo,
		ackPipeTx:     ackPipeTx,
		batchSize:     batchSize,
		flushInterval: flushInterval,
	}
}

func (h *Handler) Start(ctx context.Context) error {
	const op = "message.batchSaver.Handler.Start"
	// TODO: withFields ...
	ticker := time.NewTicker(h.flushInterval)
	defer ticker.Stop()

	batch := make([]vo.StreamMessage, 0, h.batchSize) // TODO: []entity.Message

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}

		// TODO: []entity.Message parse

		if err := h.repo.SaveBatch(ctx, batch); err != nil {
			return fmt.Errorf("batch save failed: %w", err)
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
				return err
			}

		case msg := <-h.saverReadPipe.GetReadChan():
			batch = append(batch, msg)

			if len(batch) >= h.batchSize {
				if err := flush(); err != nil {
					return err
				}
			}
		}
	}
}
