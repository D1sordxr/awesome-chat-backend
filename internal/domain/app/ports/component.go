package ports

import "context"

type Component interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
