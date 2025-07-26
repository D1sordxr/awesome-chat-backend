package wsServer

import (
	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/infrastructure/config/apps/wsServer"
	"awesome-chat/internal/infrastructure/logger"
	"awesome-chat/internal/infrastructure/ws/chathub"
	"awesome-chat/internal/presentation/httpGin/delivery/handlers/ws"
	"awesome-chat/internal/presentation/httpGin/middleware"
	"context"
	"golang.org/x/sync/errgroup"
	"time"

	ginServer "awesome-chat/internal/presentation/httpGin"
)

type App struct {
	log ports.Logger

	errChan    chan error
	components []ports.Component
}

func setupComponents(components ...ports.Component) []ports.Component {
	return components
}

func NewApp(_ context.Context) *App {
	cfg := wsServer.NewConfig()

	log := logger.NewLogger()

	wsClientManager := chathub.NewClientManager(log)

	healthHandler := ws.NewHealthHandler()

	authMid := middleware.NewAuthAsClient()
	upgradeHandler := ws.NewUpgradeHandler(wsClientManager, authMid)

	broadcastHandler := ws.NewBroadcastHandler(wsClientManager)

	server := ginServer.NewServer(
		log,
		&cfg.HTTPServer,
		upgradeHandler,
		broadcastHandler,
		healthHandler,
	)

	components := setupComponents(
		wsClientManager,
		server,
	)

	return &App{
		log:        log,
		errChan:    make(chan error),
		components: components,
	}
}

func (a *App) Run(ctx context.Context) {
	defer a.shutdown()

	errGroup, ctx := errgroup.WithContext(ctx)
	go func() { a.errChan <- errGroup.Wait() }()

	for _, component := range a.components {
		c := component
		errGroup.Go(func() error {
			return c.Start(ctx)
		})
	}

	select {
	case err := <-a.errChan:
		if err != nil {
			a.log.Error("App received an error: ", err.Error())
		}
	case <-ctx.Done():
		a.log.Info("App received a terminate signal")
	}
}

func (a *App) shutdown() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := len(a.components) - 1; i >= 0; i-- {
		c := a.components[i]
		if err := c.Shutdown(shutdownCtx); err != nil {
			a.log.Error("Error shutting down component: ", err.Error())
		}
	}
}
