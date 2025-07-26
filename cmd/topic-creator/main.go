package main

import (
	config "awesome-chat/internal/infrastructure/config/apps/topic-creator"
	"awesome-chat/internal/infrastructure/kafka"
	"log/slog"
)

func main() {
	log := slog.Default()

	log.Info("Start parsing config...")
	cfg := config.NewConfig()

	log.Info("Attempting to create topic...")
	kafka.CreateTopic(&cfg.MessageBroker)
	log.Info("Topic created successfully!")
}
