package wsserver

import (
	"fmt"

	"awesome-chat/internal/infrastructure/config/components/http"
	"awesome-chat/internal/infrastructure/config/components/redis"
	"awesome-chat/internal/infrastructure/config/loader"
)

const defaultConfigPath = "./configs/ws-server/prod.yaml"

type Config struct {
	HTTPServer http.Config  `yaml:"http"`
	Redis      redis.Config `yaml:"redis"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := loader.Read(defaultConfigPath, &cfg); err != nil {
		return nil, fmt.Errorf("ws-server config: %w", err)
	}
	return &cfg, nil
}
