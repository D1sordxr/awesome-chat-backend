package app

import (
	"awesome-chat/internal/domain/app/ports"
	"context"
	"errors"
	"fmt"
)

type LifecycleManager struct {
	components []ports.Component
}

func NewLifecycleManager(components ...ports.Component) *LifecycleManager {
	return &LifecycleManager{components: components}
}

func (lm *LifecycleManager) StartEverything(ctx context.Context) error {
	errs := make([]error, 0, len(lm.components))
	for _, component := range lm.components {
		err := component.Start(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to start component %T: %w", component, err))
		}
	}

	return errors.Join(errs...)
}

func (lm *LifecycleManager) ShutdownEverything(ctx context.Context) error {
	errs := make([]error, 0, len(lm.components))
	for _, component := range lm.components {
		err := component.Shutdown(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown component %T: %w", component, err))
		}
	}

	return errors.Join(errs...)
}
