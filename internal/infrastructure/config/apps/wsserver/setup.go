package wsserver

import (
	"fmt"

	"awesome-chat/internal/infrastructure/config/apps"
	"awesome-chat/internal/infrastructure/config/components/http"
	"awesome-chat/internal/infrastructure/config/components/redis"
	"awesome-chat/internal/infrastructure/config/env"
)

type Config struct {
	Env        apps.AppEnv
	HTTPServer http.Config
	Redis      redis.Config
}

func NewConfig() (*Config, error) {
	l := env.NewLoader()

	cfg := Config{
		HTTPServer: http.NewConfig(l),
		Redis:      redis.NewConfig(l),
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
		return nil, fmt.Errorf("ws-server config: %w", err)
	}

	return &cfg, nil
}
