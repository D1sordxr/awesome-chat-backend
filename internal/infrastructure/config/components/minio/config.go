package minio

import "awesome-chat/internal/infrastructure/config/env"

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

func NewConfig(l *env.Loader) Config {
	return Config{
		Endpoint:  l.String("MINIO_ENDPOINT"),
		AccessKey: l.String("MINIO_ACCESS_KEY"),
		SecretKey: l.String("MINIO_SECRET_KEY"),
		UseSSL:    l.BoolDefault("MINIO_USE_SSL", false),
	}
}
