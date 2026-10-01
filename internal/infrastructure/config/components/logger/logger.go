package logger

import (
	"github.com/D1sordxr/packages/log"

	"awesome-chat/internal/infrastructure/config/env"
)

func NewConfig(l *env.Loader) log.Config {
	return log.Config{
		Level:  l.StringDefault("LOG_LEVEL", "info"),
		Format: log.Format(l.StringDefault("LOG_FORMAT", string(log.FormatJSON))),
	}
}
