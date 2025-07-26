package api

import (
	"awesome-chat/internal/infrastructure/config/http"
	"awesome-chat/internal/infrastructure/config/http/wsServerApi"
	"awesome-chat/internal/infrastructure/config/jwt"
	"awesome-chat/internal/infrastructure/config/postgres"
	"awesome-chat/internal/infrastructure/config/redis"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const basicConfigPath = "./configs/api/prod.yaml"

type Config struct {
	Storage          postgres.Config    `yaml:"storage"`
	MessagePublisher redis.Config       `yaml:"message_publisher"`
	HTTPServer       http.Config        `yaml:"http"`
	WSServerAPI      wsServerApi.Config `yaml:"ws_server_api"`
	JWT              jwt.Config         `yaml:"jwt"`
}

func NewConfig() *Config {
	var cfg Config

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = basicConfigPath
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}
