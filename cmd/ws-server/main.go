package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	logging "github.com/D1sordxr/packages/log"

	"awesome-chat/internal/application/message/useCases/broadcast"
	"awesome-chat/internal/bootstrap"
	wsServer "awesome-chat/internal/infrastructure/config/apps/wsserver"
	"awesome-chat/internal/infrastructure/redis"
	"awesome-chat/internal/infrastructure/redis/stream"
	streamNames "awesome-chat/internal/infrastructure/redis/stream/names"
	"awesome-chat/internal/infrastructure/ws/chathub"
	"awesome-chat/internal/infrastructure/ws/chathub/transport"
	"awesome-chat/internal/infrastructure/ws/chathub/transport/sendMessage"
	ginServer "awesome-chat/internal/transport/httpgin"
	"awesome-chat/internal/transport/httpgin/handler/ws"
	"awesome-chat/internal/transport/httpgin/middleware"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := wsServer.NewConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err.Error())
		os.Exit(1)
	}

	log, err := logging.New(cfg.Log, os.Stdout)
	if err != nil {
		slog.Error("Failed to build logger", "error", err.Error())
		os.Exit(1)
	}

	redisConn := redis.NewConnection(&cfg.Redis)
	redisStreamPub := stream.NewPublisherImpl(redisConn, streamNames.SentMessage.String())

	wsClientStore := chathub.NewInMemoryClientStoreImpl()
	wsClientManager := chathub.NewClientManagerV2(log, wsClientStore)

	messageBroadcastWithPubUC := broadcast.NewMessageBroadcastWithPubImpl(
		log,
		redisStreamPub,
		wsClientManager,
	)
	wsSendMsgOpHandler := sendMessage.New(messageBroadcastWithPubUC)
	wsClientManager.MustSetOperationHandler(transport.NewOperationHandler(log, wsSendMsgOpHandler))

	server := ginServer.NewServer(
		log,
		&cfg.HTTPServer,
		ws.NewUpgradeHandler(wsClientManager, middleware.NewAuthAsClient()),
		ws.NewBroadcastHandler(wsClientManager),
		ws.NewHealthHandler(),
	)

	app := bootstrap.NewApp(
		log,
		redisConn,
		wsClientManager,
		server,
	)

	if err := app.Run(ctx); err != nil {
		log.Error("App exited with error", "error", err.Error())
		os.Exit(1)
	}
}
