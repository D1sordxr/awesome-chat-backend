package broadcaster

import (
	"awesome-chat/internal/application/message/useCases/broadcast"
	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/infrastructure/app"
	"awesome-chat/internal/infrastructure/config/apps/broadcaster"
	"awesome-chat/internal/infrastructure/kafka"
	"awesome-chat/internal/infrastructure/postgres"
	"awesome-chat/internal/infrastructure/postgres/executor"
	"awesome-chat/internal/infrastructure/postgres/store/user"
	"awesome-chat/internal/infrastructure/ws/chat"
	"awesome-chat/internal/presentation/http/delivery/handlers/ws"
	"awesome-chat/internal/presentation/http/middleware"
	"awesome-chat/internal/presentation/workers"
	"awesome-chat/internal/presentation/workers/kafka/handlers/message"
	"context"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"

	httpServer "awesome-chat/internal/presentation/http"
)

const workersCount = 5

type App struct {
	log        ports.Logger
	shutdowner ports.Shutdowner
	server     *httpServer.Server
	worker     *workers.Worker
}

func NewApp(ctx context.Context) *App {
	cfg := broadcaster.NewConfig()

	log := slog.Default()

	pool := postgres.NewPool(ctx, &cfg.Storage)
	txManager := executor.NewTransactionManager(pool)

	consumer := kafka.NewConsumer(&cfg.MessageBroker)

	wsChatHub := chat.NewChatHub(log)
	wsConnManager := chat.NewManager(wsChatHub)

	authStore := user.NewGetChatIDsStore(txManager)
	authMid := middleware.NewAuth(authStore)
	upgradeHandler := ws.NewUpgrade(wsConnManager, authMid)

	broadcastUC := broadcast.NewUseCase(wsChatHub)
	broadcastHandler := message.NewBroadcastMessage(
		log,
		consumer,
		broadcastUC,
	)

	server := httpServer.NewServer(
		log,
		&cfg.HTTPServer,
		upgradeHandler,
	)

	worker := workers.NewWorker(
		log,
		broadcastHandler,
	)

	shutdowner := app.NewShutdowner(
		server,
		worker,
		wsChatHub,
		pool,
	)

	return &App{
		log:        log,
		shutdowner: shutdowner,
		server:     server,
	}
}

func (a *App) Run(ctx context.Context) {
	defer a.shutdown()

	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		return a.server.Start(ctx)
	})
	//for i := 0; i < workersCount; i++ {
	//	errGroup.Go(func() error {
	//		return a.worker.Start(ctx)
	//	})
	//}

	if err := errGroup.Wait(); err != nil {
		a.log.Error("App error: " + err.Error())
	}
}

func (a *App) shutdown() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.shutdowner.ShutdownComponents(shutdownCtx); err != nil {
		a.log.Error("Failed to shutdown components: " + err.Error())
	} else {
		a.log.Info("All components shutdown successfully")
	}

	a.log.Info("App stopped gracefully")
}
