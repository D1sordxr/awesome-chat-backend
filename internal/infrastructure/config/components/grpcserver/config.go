package grpcserver

import (
	"time"

	"awesome-chat/internal/infrastructure/config/env"
)

type Config struct {
	Port            string
	ShutdownTimeout time.Duration
}

func NewConfig(l *env.Loader) Config {
	return Config{
		Port:            l.String("GRPC_PORT"),
		ShutdownTimeout: l.DurationDefault("GRPC_SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

func (c Config) Address() string {
	return ":" + c.Port
}

func (c Config) LoopbackAddress() string {
	return "localhost:" + c.Port
}
