package outbox

import (
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
	"errors"
	"fmt"

	"awesome-chat/internal/domain/core/shared/outbox/entity"
	"awesome-chat/internal/domain/core/shared/outbox/filters"

	"github.com/jackc/pgx/v5"
)

type Store struct {
	executor ports.ExecutorManager
}

func NewStore(executor ports.ExecutorManager) *Store {
	return &Store{executor: executor}
}

func (s *Store) GetMessagesForUpdate(ctx context.Context, filter filters.GetOutbox) ([]entity.Outbox, error) {
	const op = "store.outbox.GetMessagesForUpdate"

	tx, err := s.executor.GetTxExecutor(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: get transaction: %w", op, err)
	}

	query := `SELECT 
		id, 
		entity_name,
		status, 
		payload, 
		created_at
	FROM outbox
	WHERE status = $1 AND entity_name = $2
	ORDER BY created_at ASC
	LIMIT $3
	FOR UPDATE SKIP LOCKED`

	rows, err := tx.Query(ctx, query, filter.Status, filter.EntityName, filter.Limit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: query message: %w", op, err)
	}
	defer rows.Close()

	var outboxes []entity.Outbox
	for rows.Next() {
		var outbox entity.Outbox
		if err = rows.Scan(
			&outbox.OutboxID,
			&outbox.EntityName,
			&outbox.Status,
			&outbox.Payload,
		); err != nil {
			return nil, fmt.Errorf("%s: scan message: %w", op, err)
		}
		outboxes = append(outboxes, outbox)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return outboxes, nil
}

func (s *Store) UpdateStatus(ctx context.Context, status filters.SetOutboxStatus) error {
	const op = "store.outbox.UpdateStatus"

	tx, err := s.executor.GetTxExecutor(ctx)
	if err != nil {
		return fmt.Errorf("%s: get transaction: %w", op, err)
	}

	query := `UPDATE outbox SET status = $1 WHERE id = $2`
	if _, err = tx.Exec(ctx, query, status.Status, status.OutboxID); err != nil {
		return fmt.Errorf("%s: execute update: %w", op, err)
	}

	return nil
}
