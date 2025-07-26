package main

import (
	"awesome-chat/internal/bootstrap/api"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := api.NewApp(ctx)
	app.Run(ctx)
}
