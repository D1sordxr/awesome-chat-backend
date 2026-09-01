package topiccreator

import (
	"fmt"

	"awesome-chat/internal/infrastructure/config/components/kafka"
	"awesome-chat/internal/infrastructure/config/loader"
)

const defaultConfigPath = "./configs/topic-creator/prod.yaml"

type Config struct {
	MessageBroker kafka.Config `yaml:"broker"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := loader.Read(defaultConfigPath, &cfg); err != nil {
		return nil, fmt.Errorf("topic-creator config: %w", err)
	}
	return &cfg, nil
}
