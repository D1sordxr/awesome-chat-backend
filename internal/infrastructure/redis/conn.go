package redis

import (
	cfg "awesome-chat/internal/infrastructure/config/redis"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Connection struct {
	*redis.Client
}

func NewConnection(cfg *cfg.Config) *Connection {
	return &Connection{
		Client: redis.NewClient(&redis.Options{
			Addr:         cfg.GetClientAddress(),
			Password:     cfg.GetPassword(),
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			MaxRetries:   3,
		}),
	}
}

func (c *Connection) Start(ctx context.Context) error {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			return c.Client.Ping(ctx).Err()
		}
	}
}

func (c *Connection) Shutdown(_ context.Context) error {
	return c.Client.Close()
}
