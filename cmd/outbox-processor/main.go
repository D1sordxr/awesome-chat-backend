package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"awesome-chat/internal/application/outbox/useCases/process"
	"awesome-chat/internal/bootstrap"
	"awesome-chat/internal/domain/core/shared/outbox/filters"
	"awesome-chat/internal/domain/core/shared/outbox/vo"
	config "awesome-chat/internal/infrastructure/config/apps/outboxprocessor"
	"awesome-chat/internal/infrastructure/kafka"
	"awesome-chat/internal/infrastructure/logger"
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

	log := logger.NewLogger()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Error("Failed to load config", "error", err.Error())
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
