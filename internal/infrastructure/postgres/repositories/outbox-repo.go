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
			entity_name,
			status,
			payload
		) VALUES ($1, $2, $3, $4)`

	if _, err := executor.Exec(
		ctx,
		query,
		outbox.OutboxID,
		outbox.EntityName,
		outbox.Status,
		outbox.Payload,
	); err != nil {
		return err
	}

	return nil
}
