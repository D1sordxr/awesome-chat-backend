package broadcaster

import (
	"awesome-chat/internal/infrastructure/config/http"
	"awesome-chat/internal/infrastructure/config/kafka"
	"awesome-chat/internal/infrastructure/config/postgres"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const basicConfigPath = "./configs/broadcaster/prod.yaml"

type Config struct {
	Storage       postgres.Config `yaml:"storage"`
	HTTPServer    http.Config     `yaml:"http"`
	MessageBroker kafka.Config    `yaml:"broker"`
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
