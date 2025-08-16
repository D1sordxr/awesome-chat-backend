package storage

import (
	conn "awesome-chat/internal/infrastructure/redis"
	"context"
	"errors"
	"fmt"
	redisLib "github.com/redis/go-redis/v9"
	"time"
)

const defaultExpTime = time.Minute * 10

type Storage struct {
	conn   *conn.Connection
	prefix Prefix
}

func NewStorage(conn *conn.Connection, prefixes ...Prefix) *Storage {
	return &Storage{
		conn:   conn,
		prefix: NewPrefix(prefixes...),
	}
}

func (s *Storage) GetPrefix() string {
	return s.prefix.String()
}

func (s *Storage) Set(
	ctx context.Context,
	keyEndpoint string,
	value any,
	expiration ...time.Duration,
) error {
	const op = "redis.storage.Set"

	exp := defaultExpTime
	if len(expiration) > 0 {
		exp = expiration[0]
	}

	fullKey := s.prefix.WithValue(keyEndpoint)
	if err := s.conn.Set(ctx, fullKey, value, exp).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Read(ctx context.Context, keyEndpoint string) (string, error) {
	const op = "redis.storage.Read"

	fullKey := s.prefix.WithValue(keyEndpoint)
	data, err := s.conn.Get(ctx, fullKey).Result()
	if err != nil {
		if errors.Is(err, redisLib.Nil) {
			return "", fmt.Errorf("%s: %w: key %s", op, ErrKeyNotFound, fullKey)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}

func (s *Storage) Delete(ctx context.Context, keyEndpoint string) error {
	const op = "redis.storage.Delete"

	fullKey := s.prefix.WithValue(keyEndpoint)
	if err := s.conn.Del(ctx, fullKey).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
