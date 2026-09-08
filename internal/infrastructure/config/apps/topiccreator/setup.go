package topiccreator

import (
	"fmt"

	"awesome-chat/internal/infrastructure/config/apps"
	"awesome-chat/internal/infrastructure/config/components/kafka"
	"awesome-chat/internal/infrastructure/config/env"
)

type Config struct {
	Env           apps.AppEnv
	MessageBroker kafka.Config
}

func NewConfig() (*Config, error) {
	l := env.NewLoader()

	cfg := Config{
		MessageBroker: kafka.NewConfig(l),
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
		return nil, fmt.Errorf("topic-creator config: %w", err)
	}

	return &cfg, nil
}
