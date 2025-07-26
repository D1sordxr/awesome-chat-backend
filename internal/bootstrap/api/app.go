package api

import (
	chatAddMember "awesome-chat/internal/application/chat/useCases/addMember"
	chatCreate "awesome-chat/internal/application/chat/useCases/create"
	"awesome-chat/internal/application/chat/useCases/getAllMessages"
	"awesome-chat/internal/application/chat/useCases/getUserChatPreview"
	messageSendFast "awesome-chat/internal/application/message/useCases/fast"
	messageGet "awesome-chat/internal/application/message/useCases/get"
	messageSave "awesome-chat/internal/application/message/useCases/save"
	messageSend "awesome-chat/internal/application/message/useCases/send"
	"awesome-chat/internal/application/user/useCases/authJWT"
	"awesome-chat/internal/application/user/useCases/getAllUsers"
	"awesome-chat/internal/application/user/useCases/getUserChatIDs"
	"awesome-chat/internal/application/user/useCases/login"
	"awesome-chat/internal/application/user/useCases/register"
	"awesome-chat/internal/domain/app/ports"
	msgEntity "awesome-chat/internal/domain/core/message/services/entity"
	outboxEntity "awesome-chat/internal/domain/core/shared/outbox/services/entity"
	"awesome-chat/internal/infrastructure/config/apps/api"
	"awesome-chat/internal/infrastructure/jwt/user"
	"awesome-chat/internal/infrastructure/postgres"
	"awesome-chat/internal/infrastructure/postgres/executor"
	repos "awesome-chat/internal/infrastructure/postgres/repositories"
	chatStore "awesome-chat/internal/infrastructure/postgres/store/chat"
	messageStore "awesome-chat/internal/infrastructure/postgres/store/message"
	userStore "awesome-chat/internal/infrastructure/postgres/store/user"
	"awesome-chat/internal/infrastructure/redis"
	fiberHttp "awesome-chat/internal/presentation/httpFiber"
	chatHandler "awesome-chat/internal/presentation/httpFiber/delivery/handlers/chat"
	"awesome-chat/internal/presentation/httpFiber/delivery/handlers/health"
	messageHandler "awesome-chat/internal/presentation/httpFiber/delivery/handlers/message"
	userHandler "awesome-chat/internal/presentation/httpFiber/delivery/handlers/user"
	"context"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"
)

type App struct {
	log        ports.Logger
	components []ports.Component
	errChan    chan error
}

func setupComponents(components ...ports.Component) []ports.Component {
	return components
}

func NewApp(ctx context.Context) *App {
	cfg := api.NewConfig()

	log := slog.Default()

	pool := postgres.NewPool(ctx, &cfg.Storage)
	txManager := executor.NewTransactionManager(pool)

	healthHandler := new(health.Handler)

	userTokenCreator := user.NewTokenCreator(cfg.JWT.SecretKey)
	userTokenParser := user.NewTokenParser(log, cfg.JWT.SecretKey)

	userRepo := repos.NewUserRepo(txManager)
	userGetStore := userStore.NewGetStore(txManager)
	userProviderStore := userStore.NewProviderStore(txManager)
	userGetChatIDsStore := userStore.NewGetChatIDsStore(txManager)

	userRegisterUC := register.NewUserRegisterUseCase(log, userRepo)
	userLoginUC := login.NewUserLoginUseCase(log, userProviderStore, userTokenCreator)
	userAuthJwtUC := authJWT.NewUserAuthJWTUseCase(log, userTokenParser, userProviderStore)
	userGetChatIDsUC := getUserChatIDs.NewUserGetChatIDsUseCase(userGetChatIDsStore)
	userGetAllUC := getAllUsers.NewUsersGetAllUseCase(log, userGetStore)

	userHandlers := userHandler.NewUserHandler(
		userRegisterUC,
		userLoginUC,
		userAuthJwtUC,
		userGetChatIDsUC,
		userGetAllUC,
	)
	userAuthHandler := userHandler.NewAuth(log, userAuthJwtUC)

	// TODO: validateUserStore := user.NewValidateStore(txManager)
	chatCreateWithMembersStore := chatStore.NewCreateWithMembersStore(txManager)
	chatValidatorStore := chatStore.NewValidatorStore(txManager)
	chatPreviewStore := chatStore.NewGetUserChatPreviewStore(txManager)
	chatGetAllMessagesStore := chatStore.NewGetAllMessagesStore(txManager)

	chatCreateUC := chatCreate.NewChatCreateUseCase(
		txManager,
		chatCreateWithMembersStore,
		userGetStore, // TODO: change getUser -> validate userStore
	)
	chatAddMemberUC := chatAddMember.NewChatAddMemberUseCase(
		chatCreateWithMembersStore,
		chatValidatorStore,
		userGetStore,
	)
	chatPreviewUC := getUserChatPreview.NewChatGetUserChatPreviewUseCase(
		log,
		chatPreviewStore,
	)
	chatGetAllMessagesUC := getAllMessages.NewChatGetAllMessagesUseCase(
		log,
		chatGetAllMessagesStore,
	)

	chatHandlers := chatHandler.NewChatHandler(
		chatCreateUC,
		chatAddMemberUC,
		chatPreviewUC,
		chatGetAllMessagesUC,
	)

	messagePublisher := redis.NewPublisher(&cfg.MessagePublisher)
	messageGetStore := messageStore.NewGetStore(txManager)
	messageRepo := repos.NewMessageRepo(txManager)
	outboxRepo := repos.NewOutboxRepo(txManager)

	messageEntityCreator := new(msgEntity.Create)
	outboxEntityCreator := new(outboxEntity.Create)

	messageGetUC := messageGet.NewMessageGetUseCase(messageGetStore)
	messageSendUC := messageSend.NewUseCase( // TODO: rebuild
		messageEntityCreator,
		outboxEntityCreator,
		txManager,
		messageRepo,
		outboxRepo,
	)
	messageSaveUC := messageSave.NewMessageSaveUseCase(
		messageEntityCreator,
		messageRepo,
	)
	messageSendFastUC := messageSendFast.NewMessageSendFastUseCase(messagePublisher)
	messageSendSyncUC := messageSend.NewMessageSendSyncUseCase(
		messageRepo,
		messageEntityCreator,
		cfg.WSServerAPI.BroadcastURL,
	)

	messageHandlers := messageHandler.NewMessageHandler(
		messageGetUC,
		messageSaveUC,
		messageSendUC,
		messageSendFastUC,
		messageSendSyncUC,
	)

	srv := fiberHttp.NewServer(
		&cfg.HTTPServer,
		healthHandler,
		chatHandlers,
		userHandlers,
		messageHandlers,
		userAuthHandler,
	)

	components := setupComponents(
		srv,
		messageSendUC, // todo: rebuild
		pool,
	)

	return &App{
		log:        log,
		components: components,
		errChan:    make(chan error),
	}
}

func (a *App) Run(ctx context.Context) {
	defer a.shutdown()

	errGroup, ctx := errgroup.WithContext(ctx)
	go func() { a.errChan <- errGroup.Wait() }()

	for _, component := range a.components {
		func(c ports.Component) {
			errGroup.Go(func() error {
				return c.Start(ctx)
			})
		}(component)
	}

	select {
	case err := <-a.errChan:
		a.log.Error("App received an error: " + err.Error())
	case <-ctx.Done():
		a.log.Info("App received a terminate signal")
	}
}

func (a *App) shutdown() {
	a.log.Info("App shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, component := range a.components {
		if err := component.Shutdown(shutdownCtx); err != nil {
			a.log.Error("Error shutting down component: ", err.Error())
		}
	}

	a.log.Info("App shutdown successfully")
}
