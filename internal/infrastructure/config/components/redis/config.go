package redis

import "awesome-chat/internal/infrastructure/config/env"

type Config struct {
	ClientAddress string
	Password      string
	Channel       string
}

func NewConfig(l *env.Loader) Config {
	return Config{
		ClientAddress: l.String("REDIS_ADDRESS"),
		Password:      l.StringDefault("REDIS_PASSWORD", ""),
		Channel:       l.StringDefault("REDIS_CHANNEL", ""),
	}
}
