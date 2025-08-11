package main

import (
	"awesome-chat/internal/bootstrap"
	"awesome-chat/internal/domain/core/message/vo"
	"awesome-chat/internal/infrastructure/config/apps/worker"
	"awesome-chat/internal/infrastructure/logger"
	"awesome-chat/internal/infrastructure/messagePipe"
	"awesome-chat/internal/infrastructure/postgres"
	"awesome-chat/internal/infrastructure/postgres/executor"
	"awesome-chat/internal/infrastructure/postgres/repositories"
	"awesome-chat/internal/infrastructure/redis"
	"awesome-chat/internal/infrastructure/redis/stream"
	streamNames "awesome-chat/internal/infrastructure/redis/stream/names"
	"awesome-chat/internal/presentation/workers"
	"awesome-chat/internal/presentation/workers/message/handlers/acknowledger"
	"awesome-chat/internal/presentation/workers/message/handlers/batchSaver"
	"awesome-chat/internal/presentation/workers/message/handlers/streamSubscriber"
	"context"
	redisLib "github.com/redis/go-redis/v9"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := worker.NewConfig()

	log := logger.NewLogger()

	pool := postgres.NewPool(ctx, &cfg.Storage)
	txManager := executor.NewTransactionManager(pool)

	redisConn := redis.NewConnection(&cfg.StreamSubscriber)

	messageRepo := repositories.NewMessageRepo(txManager)

	messageAckPipe := messagePipe.NewMessagePipe[string]()
	messageSaverPipe := messagePipe.NewMessagePipe[vo.StreamMessage]()
	messageStreamPipe := messagePipe.NewMessagePipe[redisLib.XMessage]()
	pipeCloser := messagePipe.NewPipeCloser(
		messageAckPipe,
		messageSaverPipe,
		messageStreamPipe,
	)
	messageAckPipeTx := messagePipe.NewAckPipeTx()

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
	messageSaverHandler := batchSaver.NewHandler( // TODO batch save
		log,
		messageAckPipe,
		messageSaverPipe,
		messageRepo,
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
		pipeCloser,
		messageAckPipeTx,
		messageStreamSubscriber,
		mainWorker,
	)

	app.Run(ctx)
}
