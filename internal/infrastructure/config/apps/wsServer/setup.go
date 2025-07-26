package wsServer

import (
	"awesome-chat/internal/infrastructure/config/http"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

const basicConfigPath = "./configs/ws-server/prod.yaml"

type Config struct {
	AppEnv     string      `yaml:"app_env" env-default:"prod"`
	HTTPServer http.Config `yaml:"http"`
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
