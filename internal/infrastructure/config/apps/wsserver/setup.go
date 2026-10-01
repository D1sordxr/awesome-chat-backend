package wsserver

import (
	"fmt"

	"github.com/D1sordxr/packages/log"

	"awesome-chat/internal/infrastructure/config/apps"
	"awesome-chat/internal/infrastructure/config/components/http"
	"awesome-chat/internal/infrastructure/config/components/logger"
	"awesome-chat/internal/infrastructure/config/components/redis"
	"awesome-chat/internal/infrastructure/config/env"
)

type Config struct {
	Env        apps.AppEnv
	Log        log.Config
	HTTPServer http.Config
	Redis      redis.Config
}

func NewConfig() (*Config, error) {
	l := env.NewLoader()

	cfg := Config{
		Log:        logger.NewConfig(l),
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
