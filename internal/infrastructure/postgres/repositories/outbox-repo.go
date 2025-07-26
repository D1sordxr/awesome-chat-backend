package repositories

import (
	"awesome-chat/internal/domain/core/shared/outbox/entity"
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
)

type OutboxRepo struct {
	e ports.ExecutorManager
}

func NewOutboxRepo(e ports.ExecutorManager) *OutboxRepo {
	return &OutboxRepo{e: e}
}

func (r *OutboxRepo) Save(ctx context.Context, outbox entity.Outbox) error {
	executor := r.e.GetExecutor(ctx)
	query := `
		INSERT INTO outbox (
			id,
			payload,
			status
		) VALUES ($1, $2, $3)`

	if _, err := executor.Exec(
		ctx,
		query,
		outbox.OutboxID,
		outbox.Payload,
		outbox.Status,
	); err != nil {
		return err
	}

	return nil
}
