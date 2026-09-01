package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	chatAddMember "awesome-chat/internal/application/chat/useCases/addMember"
	chatCreate "awesome-chat/internal/application/chat/useCases/create"
	"awesome-chat/internal/application/chat/useCases/getAllMessages"
	"awesome-chat/internal/application/chat/useCases/getUserChatPreview"
	messageGet "awesome-chat/internal/application/message/useCases/get"
	"awesome-chat/internal/application/message/useCases/getForChatWithFilter"
	messageSave "awesome-chat/internal/application/message/useCases/save"
	messageSend "awesome-chat/internal/application/message/useCases/send"
	"awesome-chat/internal/application/message/useCases/sendVoice"
	"awesome-chat/internal/application/user/useCases/authJWT"
	"awesome-chat/internal/application/user/useCases/getAllUsers"
	"awesome-chat/internal/application/user/useCases/getUserChatIDs"
	"awesome-chat/internal/application/user/useCases/login"
	"awesome-chat/internal/application/user/useCases/register"
	"awesome-chat/internal/bootstrap"
	msgEntity "awesome-chat/internal/domain/core/message/services/entity"
	outboxEntity "awesome-chat/internal/domain/core/shared/outbox/services/entity"
	"awesome-chat/internal/infrastructure/config/apps/api"
	"awesome-chat/internal/infrastructure/jwt/user"
	"awesome-chat/internal/infrastructure/logger"
	"awesome-chat/internal/infrastructure/minio"
	"awesome-chat/internal/infrastructure/minio/services/bucket"
	urlSvc "awesome-chat/internal/infrastructure/minio/services/url"
	voiceStore "awesome-chat/internal/infrastructure/minio/storage/voice"
	"awesome-chat/internal/infrastructure/postgres"
	"awesome-chat/internal/infrastructure/postgres/executor"
	repos "awesome-chat/internal/infrastructure/postgres/repositories"
	chatStore "awesome-chat/internal/infrastructure/postgres/store/chat"
	messageStore "awesome-chat/internal/infrastructure/postgres/store/message"
	userStore "awesome-chat/internal/infrastructure/postgres/store/user"
	"awesome-chat/internal/infrastructure/redis"
	cacheStorage "awesome-chat/internal/infrastructure/redis/storage"
	fiberHTTP "awesome-chat/internal/presentation/httpFiber"
	chatHandler "awesome-chat/internal/presentation/httpFiber/delivery/handlers/chat"
	"awesome-chat/internal/presentation/httpFiber/delivery/handlers/health"
	messageHandler "awesome-chat/internal/presentation/httpFiber/delivery/handlers/message"
	userHandler "awesome-chat/internal/presentation/httpFiber/delivery/handlers/user"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := logger.NewLogger()

	cfg, err := api.NewConfig()
	if err != nil {
		log.Error("Failed to load config", "error", err.Error())
		os.Exit(1)
	}

	pool := postgres.NewPool(ctx, &cfg.Storage)
	txManager := executor.NewTransactionManager(pool)
	minioConn := minio.NewConnection(cfg.MinIO)
	cacheConn := redis.NewConnection(&cfg.Cache)

	userTokenCreator := user.NewTokenCreator(cfg.JWT.SecretKey)
	userTokenParser := user.NewTokenParser(log, cfg.JWT.SecretKey)

	userRepo := repos.NewUserRepo(txManager)
	userGetStore := userStore.NewGetStore(txManager)
	userProviderStore := userStore.NewProviderStore(txManager)
	userGetChatIDsStore := userStore.NewGetChatIDsStore(txManager)
	userValidatorStore := userStore.NewValidatorStore(txManager)

	userHandlers := userHandler.NewUserHandler(
		register.NewUserRegisterUseCase(log, userRepo),
		login.NewUserLoginUseCase(log, userProviderStore, userTokenCreator),
		authJWT.NewUserAuthJWTUseCase(log, userTokenParser, userProviderStore),
		getUserChatIDs.NewUserGetChatIDsUseCase(userGetChatIDsStore),
		getAllUsers.NewUsersGetAllUseCase(log, userGetStore),
	)

	chatCreateWithMembersStore := chatStore.NewCreateWithMembersStore(txManager)
	chatValidatorStore := chatStore.NewValidatorStore(txManager)

	chatHandlers := chatHandler.NewChatHandler(
		chatCreate.NewChatCreateUseCase(log, txManager, chatCreateWithMembersStore, userValidatorStore),
		chatAddMember.NewChatAddMemberUseCase(chatCreateWithMembersStore, chatValidatorStore, userValidatorStore),
		getUserChatPreview.NewChatGetUserChatPreviewUseCase(log, chatStore.NewGetUserChatPreviewStore(txManager)),
		getAllMessages.NewChatGetAllMessagesUseCase(log, chatStore.NewGetAllMessagesStore(txManager)),
	)

	messageEntityCreator := new(msgEntity.Create)
	messageRepo := repos.NewMessageRepo(txManager)

	messageSendUC := messageSend.NewUseCase( // TODO: rebuild
		messageEntityCreator,
		new(outboxEntity.Create),
		txManager,
		messageRepo,
		repos.NewOutboxRepo(txManager),
	)

	voiceBucket := bucket.Voices.String()
	bucketSvc := bucket.NewService(log, minioConn)

	messageHandlers := messageHandler.NewMessageHandler(
		messageGet.NewMessageGetUseCase(messageStore.NewGetStore(txManager)),
		messageSave.NewMessageSaveUseCase(messageEntityCreator, messageRepo),
		messageSendUC,
		messageSend.NewMessageSendSyncUseCase(messageRepo, messageEntityCreator, cfg.WSServerAPI.BroadcastURL),
		getForChatWithFilter.NewMessageGetForChatWithFilterUseCase(log, messageStore.NewGetForChatWithFilter(txManager)),
		sendVoice.NewMessageSendVoiceUseCase(
			log,
			txManager,
			messageStore.NewSaveVoiceStore(txManager),
			voiceStore.NewStorage(minioConn, voiceBucket, bucketSvc),
			urlSvc.NewURLService(
				minioConn.Client,
				voiceBucket,
				cacheStorage.NewStorage(cacheConn, cacheStorage.Voice),
			),
		),
	)

	srv := fiberHTTP.NewServer(
		&cfg.HTTPServer,
		new(health.Handler),
		chatHandlers,
		userHandlers,
		messageHandlers,
	)

	app := bootstrap.NewApp(
		log,
		pool,
		cacheConn,
		minioConn,
		messageSendUC,
		srv,
	)

	if err := app.Run(ctx); err != nil {
		log.Error("App exited with error", "error", err.Error())
		os.Exit(1)
	}
}
