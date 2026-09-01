package main

import (
	"os"

	config "awesome-chat/internal/infrastructure/config/apps/topiccreator"
	"awesome-chat/internal/infrastructure/kafka"
	"awesome-chat/internal/infrastructure/logger"
)

func main() {
	log := logger.NewLogger()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Error("Failed to load config", "error", err.Error())
		os.Exit(1)
	}

	log.Info("Attempting to create topic...", "topic", cfg.MessageBroker.Topic)
	kafka.CreateTopic(&cfg.MessageBroker)
	log.Info("Topic created successfully!")
}
