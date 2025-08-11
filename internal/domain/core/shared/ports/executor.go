package ports

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Executor defines the interface for executing database operations.
type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

type TransactionManager interface {
	BeginAndInjectTx(ctx context.Context) (context.Context, error)
	RollbackTx(ctx context.Context) error
	CommitTx(ctx context.Context) error
	ExecutorManager
}

type ExecutorManager interface {
	GetExecutor(ctx context.Context) Executor
	GetTxExecutor(ctx context.Context) (Executor, error)
	GetBatchExecutor(ctx context.Context) (Executor, error)
	GetPoolExecutor() Executor
	GetPool() *pgxpool.Pool
}
