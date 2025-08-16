package cache

import (
	"context"
	"time"
)

type Storage interface {
	GetPrefix() string
	Set(ctx context.Context, keyEndpoint string, value any, expiration ...time.Duration) error
	Read(ctx context.Context, keyEndpoint string) (string, error)
	Delete(ctx context.Context, keyEndpoint string) error
}
