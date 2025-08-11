package stream

import (
	conn "awesome-chat/internal/infrastructure/redis"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type PublisherImpl struct {
	conn       *conn.Connection
	streamName string
}

func NewPublisherImpl(
	conn *conn.Connection,
	streamName string,
) *PublisherImpl {
	return &PublisherImpl{
		conn:       conn,
		streamName: streamName,
	}
}

func (p *PublisherImpl) Publish(ctx context.Context, data map[string]any) error {
	_, err := p.conn.XAdd(ctx, &redis.XAddArgs{
		Stream: p.streamName,
		Values: data,
		MaxLen: 10000,
		Approx: true,
	}).Result()
	return err
}

func (p *PublisherImpl) PublishEvent(ctx context.Context, eventType string, payload any) error {
	data := map[string]any{
		"event":   eventType,
		"payload": payload,
		"time":    time.Now().Unix(),
	}
	return p.Publish(ctx, data)
}
