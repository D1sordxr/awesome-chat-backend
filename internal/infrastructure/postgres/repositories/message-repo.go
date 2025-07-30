package repositories

import (
	"awesome-chat/internal/domain/core/message/entity"
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
)

type MessageRepo struct {
	e ports.ExecutorManager
}

func NewMessageRepo(e ports.ExecutorManager) *MessageRepo {
	return &MessageRepo{e: e}
}

func (r *MessageRepo) Save(ctx context.Context, message entity.OldMessage) error {
	executor := r.e.GetExecutor(ctx)
	query := `
		INSERT INTO messages (
			user_id,
			chat_id,
			content
		) VALUES ($1, $2, $3)`

	if _, err := executor.Exec(ctx, query,
		message.UserID,
		message.ChatID,
		message.Content,
	); err != nil {
		return err
	}

	return nil
}
