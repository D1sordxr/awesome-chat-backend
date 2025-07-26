package postgres

import (
	"awesome-chat/internal/infrastructure/config/postgres"
	"context"

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
	return p.Pool.Ping(ctx)
}

func (p *Pool) Shutdown(_ context.Context) error {
	p.Pool.Close()
	return nil
}
