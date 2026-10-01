package api

import (
	"fmt"

	"github.com/D1sordxr/packages/log"

	"awesome-chat/internal/infrastructure/config/apps"
	"awesome-chat/internal/infrastructure/config/components/cookie"
	"awesome-chat/internal/infrastructure/config/components/grpcserver"
	"awesome-chat/internal/infrastructure/config/components/http"
	"awesome-chat/internal/infrastructure/config/components/jwt"
	"awesome-chat/internal/infrastructure/config/components/logger"
	"awesome-chat/internal/infrastructure/config/components/minio"
	"awesome-chat/internal/infrastructure/config/components/postgres"
	"awesome-chat/internal/infrastructure/config/env"
)

type Config struct {
	Env        apps.AppEnv
	Log        log.Config
	Storage    postgres.Config
	MinIO      minio.Config
	HTTPServer http.Config
	GRPCServer grpcserver.Config
	Cookie     cookie.Config
	JWT        jwt.Config

	AllowedOrigins []string
}

func NewConfig() (*Config, error) {
	l := env.NewLoader()

	cfg := Config{
		Log:        logger.NewConfig(l),
		Storage:    postgres.NewConfig(l),
		MinIO:      minio.NewConfig(l),
		HTTPServer: http.NewConfig(l),
		GRPCServer: grpcserver.NewConfig(l),
		Cookie:     cookie.NewConfig(l),
		JWT:        jwt.NewConfig(l),

		AllowedOrigins: l.StringSliceDefault("HTTP_ALLOWED_ORIGINS", ",", []string{"http://localhost:3000"}),
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
