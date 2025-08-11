package worker

import (
	"awesome-chat/internal/infrastructure/config/postgres"
	"awesome-chat/internal/infrastructure/config/redis"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

const basicConfigPath = "./configs/worker/prod.yaml"

type Config struct {
	Storage          postgres.Config `yaml:"storage"`
	StreamSubscriber redis.Config    `yaml:"stream_subscriber"`
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
