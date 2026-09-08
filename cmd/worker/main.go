package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"awesome-chat/internal/bootstrap"
	"awesome-chat/internal/domain/core/message/vo"
	workerCfg "awesome-chat/internal/infrastructure/config/apps/worker"
	"awesome-chat/internal/infrastructure/logger"
	"awesome-chat/internal/infrastructure/messagePipe"
	"awesome-chat/internal/infrastructure/postgres"
	"awesome-chat/internal/infrastructure/postgres/executor"
	"awesome-chat/internal/infrastructure/postgres/store/message"
	"awesome-chat/internal/infrastructure/redis"
	"awesome-chat/internal/infrastructure/redis/stream"
	"awesome-chat/internal/transport/worker"
	"awesome-chat/internal/transport/worker/message/acknowledger"
	"awesome-chat/internal/transport/worker/message/saver"
	"awesome-chat/internal/transport/worker/message/subscriber"

	streamNames "awesome-chat/internal/infrastructure/redis/stream/names"
	redisLib "github.com/redis/go-redis/v9"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := logger.NewLogger()

	cfg, err := workerCfg.NewConfig()
	if err != nil {
		log.Error("Failed to load config", "error", err.Error())
		os.Exit(1)
	}

	pool := postgres.NewPool(ctx, &cfg.Storage)
	txManager := executor.NewTransactionManager(pool)

	redisConn := redis.NewConnection(&cfg.StreamSubscriber)

	messageSaveFromStreamStore := message.NewSaveFromStreamStore(txManager)

	messageAckPipeTx := messagePipe.NewAckPipeTx()
	messageAckPipe := messagePipe.NewMessagePipe[string]()
	messageSaverPipe := messagePipe.NewMessagePipe[vo.StreamMessage]()
	messageStreamPipe := messagePipe.NewMessagePipe[redisLib.XMessage]()

	messageStreamSubscriber := stream.NewSubscriberImpl(
		log,
		redisConn,
		messageStreamPipe,
		streamNames.SentMessage.String(),
		streamNames.MessagesForSave.String(),
		streamNames.MessageSaverID.String(),
	)

	messageAckHandler := acknowledger.NewHandler(
		log,
		messageAckPipe,
		messageStreamSubscriber,
		messageAckPipeTx,
	)
	messageSaverHandler := saver.NewHandler(
		log,
		messageAckPipe,
		messageSaverPipe,
		messageSaveFromStreamStore,
		messageAckPipeTx,
	)
	messageReaderHandler := subscriber.NewHandler(
		log,
		messageStreamPipe,
		messageSaverPipe,
	)

	mainWorker := worker.NewWorker(
		log,
		messageAckHandler,
		messageSaverHandler,
		messageReaderHandler,
	)

	app := bootstrap.NewApp(
		log,
		pool,
		redisConn,
		messageAckPipeTx,
		//pipeCloser,
		messageStreamSubscriber,
		mainWorker,
	)

	if err := app.Run(ctx); err != nil {
		log.Error("App exited with error", "error", err.Error())
		os.Exit(1)
	}
}
