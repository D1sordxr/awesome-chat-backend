package executor

import (
	"awesome-chat/internal/domain/core/shared/ports"
	"awesome-chat/internal/infrastructure/postgres"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgx/v5"
)

// txKey and batchKey are used as keys for storing transactions and batches in the context.
type (
	txKey    struct{}
	batchKey struct{}
)

// TransactionManager of the Executor interface.
// It is used to manage transactions and batches in the context and delegate queries to the appropriate executor.
type TransactionManager struct {
	*postgres.Pool
}

// NewTransactionManager creates a new TransactionManager instance with the given Postgres connection pool.
func NewTransactionManager(pool *postgres.Pool) *TransactionManager {
	return &TransactionManager{Pool: pool}
}

func (m *TransactionManager) BeginAndInjectTx(ctx context.Context) (context.Context, error) {
	tx, err := m.Pool.Begin(ctx)
	if err != nil {
		return ctx, err
	}

	ctx = m.injectTx(ctx, tx)

	return ctx, nil
}

func (m *TransactionManager) CommitTx(ctx context.Context) error {
	tx, ok := m.extractTx(ctx)
	if !ok {
		return errors.New("transaction not found in transaction manager")
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (m *TransactionManager) RollbackTx(ctx context.Context) error {
	tx, ok := m.extractTx(ctx)
	if !ok {
		return errors.New("transaction not found in transaction manager")
	}

	if err := tx.Rollback(ctx); err != nil {
		return err
	}

	return nil
}

// injectTx stores a transaction in the context for later retrieval.
func (m *TransactionManager) injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx retrieves a transaction from the context.
// It returns the transaction and a boolean indicating whether a transaction was found.
func (m *TransactionManager) extractTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// NewBatch creates a new BatchExecutor for queueing batch queries.
func (m *TransactionManager) NewBatch() *BatchExecutor {
	return &BatchExecutor{Batch: &pgx.Batch{}}
}

// InjectBatch stores a batch in the context for later retrieval.
func (m *TransactionManager) InjectBatch(ctx context.Context, batch *BatchExecutor) context.Context {
	return context.WithValue(ctx, batchKey{}, batch)
}

// ExtractBatch retrieves a batch from the context, if it exists.
func (m *TransactionManager) ExtractBatch(ctx context.Context) (*BatchExecutor, bool) {
	batch, ok := ctx.Value(batchKey{}).(*BatchExecutor)
	return batch, ok
}

// GetExecutor returns the appropriate executor based on the context.
// If a batch is present in the context, it returns the batch executor.
// If a transaction is present, it returns the transaction.
// Otherwise, it returns a PoolExecutor, which wraps the connection pool.
func (m *TransactionManager) GetExecutor(ctx context.Context) ports.Executor {
	if batch, ok := m.ExtractBatch(ctx); ok {
		return batch
	}

	if tx, ok := m.extractTx(ctx); ok {
		return tx
	}

	return &PoolExecutor{Pool: m.Pool}
}

// GetPoolExecutor returns a PoolExecutor that wraps the connection pool.
// It can be used to execute queries outside a transaction or batch.
// Prefer using GetExecutor instead of this method.
func (m *TransactionManager) GetPoolExecutor() ports.Executor {
	return &PoolExecutor{Pool: m.Pool}
}

// GetTxExecutor returns the transaction executor from the context.
// Prefer using GetExecutor instead of this method.
func (m *TransactionManager) GetTxExecutor(ctx context.Context) (ports.Executor, error) {
	tx, ok := m.extractTx(ctx)
	if ok {
		return tx, nil
	}

	return nil, errors.New("no transaction found in context")
}

// GetBatchExecutor returns the batch executor from the context.
// Prefer using GetExecutor instead of this method.
func (m *TransactionManager) GetBatchExecutor(ctx context.Context) (ports.Executor, error) {
	batch, ok := m.ExtractBatch(ctx)
	if ok {
		return batch, nil
	}

	return nil, errors.New("no batch found in context")
}

func (m *TransactionManager) GetPool() *pgxpool.Pool {
	return m.Pool.Pool
}
