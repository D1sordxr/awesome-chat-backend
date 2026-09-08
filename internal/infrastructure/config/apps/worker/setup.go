package worker

import (
	"fmt"

	"awesome-chat/internal/infrastructure/config/apps"
	"awesome-chat/internal/infrastructure/config/components/postgres"
	"awesome-chat/internal/infrastructure/config/components/redis"
	"awesome-chat/internal/infrastructure/config/env"
)

type Config struct {
	Env              apps.AppEnv
	Storage          postgres.Config
	StreamSubscriber redis.Config
}

func NewConfig() (*Config, error) {
	l := env.NewLoader()

	cfg := Config{
		Storage:          postgres.NewConfig(l),
		StreamSubscriber: redis.NewConfig(l),
	}

	l.Custom("APP_ENV", func(raw string) error {
		appEnv, err := apps.ParseAppEnv(raw)
		if err != nil {
			return err
		}
		cfg.Env = appEnv

		return nil
	})

	if err := l.Err(); err != nil {
		return nil, fmt.Errorf("worker config: %w", err)
	}

	return &cfg, nil
}
