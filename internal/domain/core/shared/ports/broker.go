package ports

import (
	"awesome-chat/internal/domain/core/shared/broker/entity"
	"context"
)

type Producer interface {
	Publish(ctx context.Context, message entity.Message) error
}

type Consumer interface {
	Receive(ctx context.Context) (entity.Message, error)
}
