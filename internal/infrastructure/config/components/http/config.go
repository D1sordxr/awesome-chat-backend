package http

import (
	"time"

	"awesome-chat/internal/infrastructure/config/env"
)

type Config struct {
	Port        string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func NewConfig(l *env.Loader) Config {
	return Config{
		Port:        l.String("HTTP_PORT"),
		Timeout:     l.DurationDefault("HTTP_TIMEOUT", 4*time.Second),
		IdleTimeout: l.DurationDefault("HTTP_IDLE_TIMEOUT", 30*time.Second),
	}
}
