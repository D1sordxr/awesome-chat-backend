package main

import (
	"context"
	"os/signal"
	"syscall"

	"awesome-chat/internal/bootstrap"
	"awesome-chat/internal/domain/core/message/vo"
	"awesome-chat/internal/infrastructure/config/apps/worker"
	"awesome-chat/internal/infrastructure/logger"
	"awesome-chat/internal/infrastructure/messagePipe"
	"awesome-chat/internal/infrastructure/postgres"
	"awesome-chat/internal/infrastructure/postgres/executor"
	"awesome-chat/internal/infrastructure/postgres/store/message"
	"awesome-chat/internal/infrastructure/redis"
	"awesome-chat/internal/infrastructure/redis/stream"
	"awesome-chat/internal/presentation/workers"
	"awesome-chat/internal/presentation/workers/message/handlers/acknowledger"
	"awesome-chat/internal/presentation/workers/message/handlers/batchSaver"
	"awesome-chat/internal/presentation/workers/message/handlers/streamSubscriber"

	streamNames "awesome-chat/internal/infrastructure/redis/stream/names"
	redisLib "github.com/redis/go-redis/v9"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := worker.NewConfig()

	log := logger.NewLogger()

	pool := postgres.NewPool(ctx, &cfg.Storage)
	txManager := executor.NewTransactionManager(pool)

	redisConn := redis.NewConnection(&cfg.StreamSubscriber)

	messageSaveFromStreamStore := message.NewSaveFromStreamStore(txManager)

	messageAckPipeTx := messagePipe.NewAckPipeTx()
	messageAckPipe := messagePipe.NewMessagePipe[string]()
	messageSaverPipe := messagePipe.NewMessagePipe[vo.StreamMessage]()
	messageStreamPipe := messagePipe.NewMessagePipe[redisLib.XMessage]()
	pipeCloser := messagePipe.NewPipeCloser(
		messageAckPipe,
		messageSaverPipe,
		messageStreamPipe,
	)

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
	messageSaverHandler := batchSaver.NewHandler(
		log,
		messageAckPipe,
		messageSaverPipe,
		messageSaveFromStreamStore,
		messageAckPipeTx,
	)
	messageReaderHandler := streamSubscriber.NewHandler(
		log,
		messageStreamPipe,
		messageSaverPipe,
	)

	mainWorker := workers.NewWorker(
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
		pipeCloser,
		messageStreamSubscriber,
		mainWorker,
	)

	app.Run(ctx)
}
