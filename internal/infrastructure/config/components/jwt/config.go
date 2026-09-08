package jwt

import "awesome-chat/internal/infrastructure/config/env"

type Config struct {
	SecretKey string
}

func NewConfig(l *env.Loader) Config {
	return Config{
		SecretKey: l.String("JWT_SECRET_KEY"),
	}
}
