package postgres

import (
	"awesome-chat/internal/infrastructure/config/postgres"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
}

func NewPool(ctx context.Context, config *postgres.Config) *Pool {
	pool, err := pgxpool.New(ctx, config.ConnectionString())
	if err != nil {
		panic(err)
	}

	return &Pool{Pool: pool}
}

func (p *Pool) Start(ctx context.Context) error {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			return p.Pool.Ping(ctx)
		}
	}
}

func (p *Pool) Shutdown(_ context.Context) error {
	p.Pool.Close()
	return nil
}
