package main

import (
	"context"
	"os/signal"
	"syscall"

	"awesome-chat/internal/bootstrap/broadcaster"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := broadcaster.NewApp(ctx)
	app.Run(ctx)
}
