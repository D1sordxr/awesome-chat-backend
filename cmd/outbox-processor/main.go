package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	logging "github.com/D1sordxr/packages/log"

	"awesome-chat/internal/application/outbox/useCases/process"
	"awesome-chat/internal/bootstrap"
	"awesome-chat/internal/domain/core/shared/outbox/filters"
	"awesome-chat/internal/domain/core/shared/outbox/vo"
	config "awesome-chat/internal/infrastructure/config/apps/outboxprocessor"
	"awesome-chat/internal/infrastructure/kafka"
	"awesome-chat/internal/infrastructure/postgres"
	"awesome-chat/internal/infrastructure/postgres/executor"
	outboxStores "awesome-chat/internal/infrastructure/postgres/store/outbox"
	"awesome-chat/internal/transport/worker"
	"awesome-chat/internal/transport/worker/outbox"
)

const outboxBatchLimit = 10

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err.Error())
		os.Exit(1)
	}

	log, err := logging.New(cfg.Log, os.Stdout)
	if err != nil {
		slog.Error("Failed to build logger", "error", err.Error())
		os.Exit(1)
	}

	pool := postgres.NewPool(ctx, &cfg.Storage)
	txManager := executor.NewTransactionManager(pool)
	producer := kafka.NewProducer(&cfg.MessageBroker)

	processUC := process.NewUseCase(
		outboxStores.NewStore(txManager),
		txManager,
		producer,
	)

	processHandler := outbox.NewHandler(log, processUC, filters.GetOutbox{
		EntityName: vo.MessageEntity,
		Status:     vo.StatusPending,
		Limit:      outboxBatchLimit,
	})

	app := bootstrap.NewApp(
		log,
		pool,
		producer,
		worker.NewWorker(log, processHandler),
	)

	if err := app.Run(ctx); err != nil {
		log.Error("App exited with error", "error", err.Error())
		os.Exit(1)
	}
}
