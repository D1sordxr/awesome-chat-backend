package ports

import (
	"context"
)

type Subscriber interface {
	GetSubChannel(ctx context.Context) (<-chan []byte, error)
}

type StreamSubscriber interface {
	Ack(ctx context.Context, msgID string) error
}

type Acknowledger interface {
	Ack(ctx context.Context, msgID string) error
}

type StreamPublisher interface {
	Publish(ctx context.Context, data map[string]any) error
}
