package outboxprocessor

import (
	"fmt"

	"awesome-chat/internal/infrastructure/config/components/kafka"
	"awesome-chat/internal/infrastructure/config/components/postgres"
	"awesome-chat/internal/infrastructure/config/loader"
)

const defaultConfigPath = "./configs/outbox-processor/prod.yaml"

type Config struct {
	Storage       postgres.Config `yaml:"storage"`
	MessageBroker kafka.Config    `yaml:"broker"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := loader.Read(defaultConfigPath, &cfg); err != nil {
		return nil, fmt.Errorf("outbox-processor config: %w", err)
	}
	return &cfg, nil
}
