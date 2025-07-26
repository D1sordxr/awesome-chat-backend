package executor

import (
	"awesome-chat/internal/infrastructure/postgres"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// PoolExecutor is a wrapper around postgres.Pool that implements the Executor interface.
type PoolExecutor struct {
	Pool *postgres.Pool
}

// Exec executes a query without returning any rows.
func (p *PoolExecutor) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return p.Pool.Exec(ctx, sql, args...)
}

// Query executes a query that returns rows.
func (p *PoolExecutor) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.Pool.Query(ctx, sql, args...)
}

// QueryRow executes a query that is expected to return at most one row.
func (p *PoolExecutor) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.Pool.QueryRow(ctx, sql, args...)
}

// SendBatch sends a batch of queries and returns the batch results.
func (p *PoolExecutor) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return p.Pool.SendBatch(ctx, b)
}

// CopyFrom performs a bulk copy operation.
func (p *PoolExecutor) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return p.Pool.CopyFrom(ctx, tableName, columnNames, rowSrc)
}
