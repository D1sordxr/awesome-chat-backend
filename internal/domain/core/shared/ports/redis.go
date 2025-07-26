package ports

import "context"

type Subscriber interface {
	GetSubChannel(ctx context.Context) (<-chan []byte, error)
}

type Publisher interface {
	Publish(ctx context.Context, payload []byte) error
}
