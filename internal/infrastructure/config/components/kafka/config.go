package kafka

import "awesome-chat/internal/infrastructure/config/env"

type Config struct {
	Brokers []string
	Topic   string
}

func NewConfig(l *env.Loader) Config {
	return Config{
		Brokers: l.StringSlice("KAFKA_BROKERS", ","),
		Topic:   l.String("KAFKA_TOPIC"),
	}
}
