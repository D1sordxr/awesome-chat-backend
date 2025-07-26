package executor

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// BatchExecutor is a struct that implements the Executor interface for queueing batch queries.
// It allows multiple queries to be queued and executed together as a single batch.
type BatchExecutor struct {
	Batch *pgx.Batch
}

// Exec queues a SQL query with the given arguments for batch execution.
// It does not execute the query immediately but adds it to the batch.
func (b *BatchExecutor) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	b.Batch.Queue(sql, arguments...)
	return pgconn.CommandTag{}, nil
}

// Query queues a SQL query for batch execution.
// It does not execute the query immediately but adds it to the batch.
func (b *BatchExecutor) Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error) {
	b.Batch.Queue(sql, optionsAndArgs...)
	return nil, nil
}

// QueryRow queues a SQL query for batch execution.
// It does not execute the query immediately but adds it to the batch.
func (b *BatchExecutor) QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row {
	b.Batch.Queue(sql, optionsAndArgs...)
	return nil
}

// SendBatch is a placeholder method that does nothing in the BatchExecutor.
// Batch queries are executed when the batch is sent using the underlying executor.
func (b *BatchExecutor) SendBatch(ctx context.Context, batch *pgx.Batch) pgx.BatchResults {
	return nil
}

// CopyFrom is not supported in the BatchExecutor.
// It returns an error indicating that bulk copy operations are not supported in batch mode.
func (b *BatchExecutor) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, errors.New("not supported")
}
