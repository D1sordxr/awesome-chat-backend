package main

import (
	"log/slog"
	"os"

	logging "github.com/D1sordxr/packages/log"

	config "awesome-chat/internal/infrastructure/config/apps/topiccreator"
	"awesome-chat/internal/infrastructure/kafka"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err.Error())
		os.Exit(1)
	}

	log, err := logging.New(cfg.Log, os.Stdout)
	if err != nil {
		slog.Error("Failed to build logger", "error", err.Error())
		os.Exit(1)
	}

	log.Info("Attempting to create topic...", "topic", cfg.MessageBroker.Topic)
	kafka.CreateTopic(&cfg.MessageBroker)
	log.Info("Topic created successfully!")
}
