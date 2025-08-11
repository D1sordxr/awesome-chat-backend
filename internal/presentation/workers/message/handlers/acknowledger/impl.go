package acknowledger

import (
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/message/ports/worker"
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
)

type Handler struct {
	log          appPorts.Logger
	readPipe     worker.MessagePipe[string]
	acknowledger ports.Acknowledger
	ackPipeTx    worker.AckPipeTx
}

func NewHandler(
	log appPorts.Logger,
	readPipe worker.MessagePipe[string],
	acknowledger ports.Acknowledger,
	ackPipeTx worker.AckPipeTx,
) *Handler {
	return &Handler{
		log:          log,
		readPipe:     readPipe,
		acknowledger: acknowledger,
		ackPipeTx:    ackPipeTx,
	}
}

func (h *Handler) Start(ctx context.Context) error {
	const op = "message.acknowledger.Handler.Start"
	withFields := func(args ...any) []any {
		return append([]any{"operation", op}, args...)
	}
	readChan := h.readPipe.GetReadChan()

	h.log.Info("Starting acknowledger...", withFields()...)

	go func() {
		for {
			select {
			case id, ok := <-readChan:
				if !ok {
					return
				}
				if err := func(id string) error {
					defer h.ackPipeTx.Confirm()
					return h.acknowledger.Ack(ctx, id)
				}(id); err != nil {
					h.log.Error("acknowledgement error", withFields("id", id, "error", err.Error())...)
				}
			}
		}

	}()

	<-ctx.Done()
	h.ackPipeTx.Wait()

	return nil
}
