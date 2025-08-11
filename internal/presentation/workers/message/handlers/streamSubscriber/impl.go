package streamSubscriber

import (
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/message/ports/worker"
	"awesome-chat/internal/domain/core/message/vo"
	"context"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	log        appPorts.Logger
	subscriber worker.MessagePipe[redis.XMessage]
	saverPipe  worker.MessagePipe[vo.StreamMessage]
}

func NewHandler(
	log appPorts.Logger,
	streamPipe worker.MessagePipe[redis.XMessage],
	saver worker.MessagePipe[vo.StreamMessage],
) *Handler {
	return &Handler{
		log:        log,
		subscriber: streamPipe,
		saverPipe:  saver,
	}
}

func (h *Handler) Start(ctx context.Context) error {
	const op = "message.streamSubscriber.Handler.Start"
	logFields := func(args ...any) []any {
		return append([]any{"operation", op}, args...)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-h.subscriber.GetReadChan():
			if !ok {
				h.log.Info("Message channel closed", logFields()...)
				return nil
			}

			parsedMsg, err := vo.ParseStreamMessage(msg.ID, msg.Values)
			if err != nil {
				h.log.Error("Parse error", logFields("id", msg.ID, "error", err.Error())...)
				continue
			}

			select {
			case h.saverPipe.GetWriteChan() <- parsedMsg:
			case <-ctx.Done():
				return nil
			default:
				h.log.Warn("Saver queue overflow, message dropped", logFields("id", msg.ID)...)
			}
		}
	}
}
