package ports

import "context"

type Component interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type AppLifecycleManager interface {
	StartEverything(ctx context.Context) error
	ShutdownEverything(ctx context.Context) error
}

type AppComponent interface {
	Shutdown(ctx context.Context) error
}

type Shutdowner interface {
	ShutdownComponents(ctx context.Context) error
}
