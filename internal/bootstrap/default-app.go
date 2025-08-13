package bootstrap

import (
	"awesome-chat/internal/domain/app/ports"
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"time"
)

type App struct {
	log        ports.Logger
	components []ports.Component
}

func NewApp(
	log ports.Logger,
	components ...ports.Component,
) *App {
	return &App{
		log:        log,
		components: components,
	}
}

func (a *App) Run(ctx context.Context) {
	defer a.shutdown()

	errChan := make(chan error)
	errGroup, ctx := errgroup.WithContext(ctx)
	go func() { errChan <- errGroup.Wait() }()

	for _, component := range a.components {
		func(c ports.Component) {
			errGroup.Go(func() error {
				return c.Start(ctx)
			})
		}(component)
	}

	select {
	case err := <-errChan:
		a.log.Error("App received an error", "error", err.Error())
	case <-ctx.Done():
		a.log.Info("App received a terminate signal")
	}
}

func (a *App) shutdown() {
	a.log.Info("App shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	errs := make([]error, 0, len(a.components))
	for i := len(a.components) - 1; i >= 0; i-- {
		a.log.Info("Shutting down component", "idx", i)
		if err := a.components[i].Shutdown(shutdownCtx); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		a.log.Info("App successfully shutdown")
	} else {
		a.log.Error(
			"App shutdown with errors",
			"errors", errors.Join(errs...).Error(),
		)
	}
}
