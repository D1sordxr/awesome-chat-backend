package api

import (
	"fmt"

	"awesome-chat/internal/infrastructure/config/apps"
	"awesome-chat/internal/infrastructure/config/components/http"
	"awesome-chat/internal/infrastructure/config/components/jwt"
	"awesome-chat/internal/infrastructure/config/components/minio"
	"awesome-chat/internal/infrastructure/config/components/postgres"
	"awesome-chat/internal/infrastructure/config/components/redis"
	"awesome-chat/internal/infrastructure/config/components/wsserverapi"
	"awesome-chat/internal/infrastructure/config/env"
)

type Config struct {
	Env         apps.AppEnv
	Storage     postgres.Config
	Cache       redis.Config
	MinIO       minio.Config
	HTTPServer  http.Config
	WSServerAPI wsserverapi.Config
	JWT         jwt.Config
}

func NewConfig() (*Config, error) {
	l := env.NewLoader()

	cfg := Config{
		Storage:     postgres.NewConfig(l),
		Cache:       redis.NewConfig(l),
		MinIO:       minio.NewConfig(l),
		HTTPServer:  http.NewConfig(l),
		WSServerAPI: wsserverapi.NewConfig(l),
		JWT:         jwt.NewConfig(l),
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
		return nil, fmt.Errorf("api config: %w", err)
	}

	return &cfg, nil
}
