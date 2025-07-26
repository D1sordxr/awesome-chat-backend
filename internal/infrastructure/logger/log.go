package logger

import (
	"awesome-chat/internal/infrastructure/logger/handlers/slogpretty"
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
