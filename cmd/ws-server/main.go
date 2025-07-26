package main

import (
	"context"
	"os/signal"
	"syscall"

	"awesome-chat/internal/bootstrap/wsServer"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := wsServer.NewApp(ctx)
	app.Run(ctx)
}
