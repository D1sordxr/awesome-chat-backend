package pubsub

import (
	cfg "awesome-chat/internal/infrastructure/config/redis"
	"context"
	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	channel string
	client  *redis.Client
}

func NewPublisher(cfg *cfg.Config) *Publisher {
	return &Publisher{
		channel: cfg.GetChannel(),
		client: redis.NewClient(&redis.Options{
			Addr: cfg.GetClientAddress(),
		}),
	}
}

func (p *Publisher) Publish(ctx context.Context, payload []byte) error {
	if err := p.client.Publish(ctx, p.channel, payload).Err(); err != nil {
		return err
	}

	return nil
}

func (p *Publisher) Start(ctx context.Context) error {
	return p.client.Ping(ctx).Err()
}

func (p *Publisher) Shutdown(_ context.Context) error {
	return p.client.Close()
}
