package outbox

import (
	"awesome-chat/internal/application/outbox/useCases/process"
	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/shared/outbox/filters"
	"awesome-chat/internal/domain/core/shared/outbox/vo"
	"awesome-chat/internal/infrastructure/app"
	config "awesome-chat/internal/infrastructure/config/apps/outbox-processor"
	"awesome-chat/internal/infrastructure/kafka"
	"awesome-chat/internal/infrastructure/postgres"
	"awesome-chat/internal/infrastructure/postgres/executor"
	outboxStores "awesome-chat/internal/infrastructure/postgres/store/outbox"
	"awesome-chat/internal/presentation/workers"
	"awesome-chat/internal/presentation/workers/outbox/handlers"
	"context"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"
)

const workersCount = 5

type App struct {
	log        ports.Logger
	shutdowner ports.Shutdowner
	worker     *workers.Worker
}

func NewApp(ctx context.Context) *App {
	cfg := config.NewConfig()

	log := slog.Default()

	pool := postgres.NewPool(ctx, &cfg.Storage)
	txManager := executor.NewTransactionManager(pool)

	producer := kafka.NewProducer(&cfg.MessageBroker)

	outboxStore := outboxStores.NewStore(txManager)
	processUC := process.NewUseCase(
		outboxStore,
		txManager,
		producer,
	)
	processHandler := handlers.NewHandler(log, processUC, filters.GetOutbox{
		EntityName: vo.MessageEntity,
		Status:     vo.StatusPending,
		Limit:      10,
	})

	worker := workers.NewWorker(
		log,
		processHandler,
	)

	shutdowner := app.NewShutdowner(
		worker,
		producer,
		pool,
	)

	return &App{
		log:        log,
		shutdowner: shutdowner,
		worker:     worker,
	}
}

func (a *App) Run(ctx context.Context) {
	defer a.shutdown()
	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < workersCount; i++ {
		g.Go(func() error {
			return a.worker.Start(ctx)
		})
	}

	if err := g.Wait(); err != nil {
		a.log.Error("Worker error: " + err.Error())
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
