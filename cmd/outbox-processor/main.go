package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	processor "awesome-chat/internal/bootstrap/outboxProcessor"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := processor.NewApp(ctx)
	app.Run(ctx)
}
