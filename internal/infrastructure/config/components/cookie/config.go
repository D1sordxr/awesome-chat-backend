package cookie

import (
	"time"

	"awesome-chat/internal/infrastructure/config/env"
)

type Config struct {
	Name   string
	Path   string
	Domain string
	TTL    time.Duration
	Secure bool
}

func NewConfig(l *env.Loader) Config {
	return Config{
		Name:   l.StringDefault("AUTH_COOKIE_NAME", "jwt"),
		Path:   l.StringDefault("AUTH_COOKIE_PATH", "/"),
		Domain: l.StringDefault("AUTH_COOKIE_DOMAIN", ""),
		TTL:    l.DurationDefault("AUTH_COOKIE_TTL", 48*time.Hour),
		Secure: l.BoolDefault("AUTH_COOKIE_SECURE", false),
	}
}
