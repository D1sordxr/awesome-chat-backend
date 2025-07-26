package app

import (
	"awesome-chat/internal/domain/app/ports"
	"context"
	"errors"
	"fmt"
)

type Shutdowner struct {
	components []ports.AppComponent
}

func NewShutdowner(components ...ports.AppComponent) *Shutdowner {
	return &Shutdowner{
		components: components,
	}
}

func (s *Shutdowner) ShutdownComponents(ctx context.Context) error {
	errs := make([]error, 0, len(s.components))
	for _, component := range s.components {
		err := component.Shutdown(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown component %T: %w", component, err))
		}
	}

	return errors.Join(errs...)
}
