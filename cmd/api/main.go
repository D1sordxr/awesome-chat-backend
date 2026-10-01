package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	logging "github.com/D1sordxr/packages/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	chatApp "awesome-chat/internal/application/chat"
	messageApp "awesome-chat/internal/application/message"
	userApp "awesome-chat/internal/application/user"
	"awesome-chat/internal/bootstrap"
	"awesome-chat/internal/infrastructure/config/apps/api"
	jwtUser "awesome-chat/internal/infrastructure/jwt/user"
	"awesome-chat/internal/infrastructure/minio"
	"awesome-chat/internal/infrastructure/minio/services/bucket"
	"awesome-chat/internal/infrastructure/minio/services/upload"
	"awesome-chat/internal/infrastructure/postgres"
	"awesome-chat/internal/infrastructure/postgres/executor"
	repos "awesome-chat/internal/infrastructure/postgres/repositories"
	chatStore "awesome-chat/internal/infrastructure/postgres/store/chat"
	messageStore "awesome-chat/internal/infrastructure/postgres/store/message"
	userStore "awesome-chat/internal/infrastructure/postgres/store/user"
	"awesome-chat/internal/transport/gateway"
	grpcTransport "awesome-chat/internal/transport/grpc"
	chatHandler "awesome-chat/internal/transport/grpc/handler/chat"
	messageHandler "awesome-chat/internal/transport/grpc/handler/message"
	userHandler "awesome-chat/internal/transport/grpc/handler/user"
	"awesome-chat/internal/transport/grpc/interceptor"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := api.NewConfig()
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
	minioConn := minio.NewConnection(cfg.MinIO)

	userUseCase := userApp.NewUseCaseWithTracing(userApp.NewUseCase(
		repos.NewUserRepo(txManager),
		userStore.NewProviderStore(txManager),
		userStore.NewGetStore(txManager),
		userStore.NewGetChatIDsStore(txManager),
		jwtUser.NewTokenCreator(cfg.JWT.SecretKey),
	))

	chatUseCase := chatApp.NewUseCaseWithTracing(chatApp.NewUseCase(
		chatStore.NewCreateWithMembersStore(txManager),
		chatStore.NewValidatorStore(txManager),
		userStore.NewValidatorStore(txManager),
		chatStore.NewGetUserChatPreviewStore(txManager),
		txManager,
	))

	messageUseCase := messageApp.NewUseCaseWithTracing(messageApp.NewUseCase(
		repos.NewMessageRepo(txManager),
		messageStore.NewSaveVoiceStore(txManager),
		repos.NewOutboxRepo(txManager),
		messageStore.NewGetForChatWithFilter(txManager),
		upload.NewService(minioConn.Client, bucket.Voices.String()),
		chatStore.NewValidatorStore(txManager),
		txManager,
	))

	users := userHandler.NewHandlerWithTracing(userHandler.NewHandler(userUseCase))
	chats := chatHandler.NewHandlerWithTracing(chatHandler.NewHandler(chatUseCase))
	messages := messageHandler.NewHandlerWithTracing(messageHandler.NewHandler(messageUseCase))

	grpcLogger := interceptor.NewLogger(log)
	grpcErrors := interceptor.NewError(log)
	auth := interceptor.NewAuth(
		jwtUser.NewTokenParser(log, cfg.JWT.SecretKey),
		interceptor.PublicMethods...,
	)

	grpcServer := grpcTransport.NewServer(
		log,
		cfg.GRPCServer,
		[]grpc.UnaryServerInterceptor{grpcLogger.Unary(), grpcErrors.Unary(), auth.Unary()},
		func(server grpc.ServiceRegistrar) { v1.RegisterUserServiceServer(server, users) },
		func(server grpc.ServiceRegistrar) { v1.RegisterChatServiceServer(server, chats) },
		func(server grpc.ServiceRegistrar) { v1.RegisterMessageServiceServer(server, messages) },
	)

	conn, err := grpc.NewClient(
		cfg.GRPCServer.LoopbackAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to dial gRPC server", "error", err.Error())
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()

	httpGateway, err := gateway.NewServer(
		ctx,
		log,
		&cfg.HTTPServer,
		gateway.Options{
			AllowedOrigins: cfg.AllowedOrigins,
			Cookie:         cfg.Cookie,
		},
		conn,
		v1.RegisterUserServiceHandler,
		v1.RegisterChatServiceHandler,
		v1.RegisterMessageServiceHandler,
	)
	if err != nil {
		log.Error("Failed to build HTTP gateway", "error", err.Error())
		os.Exit(1)
	}

	app := bootstrap.NewApp(
		log,
		pool,
		minioConn,
		grpcServer,
		httpGateway,
	)

	if err = app.Run(ctx); err != nil {
		log.Error("App exited with error", "error", err.Error())
		os.Exit(1)
	}
}
