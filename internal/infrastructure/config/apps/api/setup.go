package api

import (
	"awesome-chat/internal/infrastructure/config/apps"
	"fmt"
	"os"
	"strconv"

	"awesome-chat/internal/infrastructure/config/components/http"
	"awesome-chat/internal/infrastructure/config/components/jwt"
	"awesome-chat/internal/infrastructure/config/components/minio"
	"awesome-chat/internal/infrastructure/config/components/postgres"
	"awesome-chat/internal/infrastructure/config/components/redis"
	"awesome-chat/internal/infrastructure/config/components/wsserverapi"
	"awesome-chat/internal/infrastructure/config/loader"
)

const defaultConfigPath = "./configs/api/prod.yaml"

type Config struct {
	Env         apps.AppEnv
	Storage     postgres.Config    `yaml:"storage"`
	Cache       redis.Config       `yaml:"cache"`
	MinIO       minio.Config       `yaml:"minio"`
	HTTPServer  http.Config        `yaml:"http"`
	WSServerAPI wsserverapi.Config `yaml:"ws_server_api"`
	JWT         jwt.Config         `yaml:"jwt"`
}

func NewConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	strconv.ParseInt(env, 10, 8)
	var cfg Config
	if err := loader.Read(defaultConfigPath, &cfg); err != nil {
		return nil, fmt.Errorf("api config: %w", err)
	}
	return &cfg, nil
}
