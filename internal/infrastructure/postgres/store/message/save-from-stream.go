package message

import (
	"awesome-chat/internal/domain/core/message/vo"
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type SaveFromStreamStore struct {
	executor ports.ExecutorManager
}

func NewSaveFromStreamStore(executor ports.ExecutorManager) *SaveFromStreamStore {
	return &SaveFromStreamStore{executor: executor}
}

func (s *SaveFromStreamStore) SaveBatch(ctx context.Context, messages []vo.StreamMessage) error {
	const op = "message.SaveFromStreamStore.SaveBatch"

	if len(messages) == 0 {
		return nil
	}

	tx, err := s.executor.GetPool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin transaction failed: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"messages"},
		[]string{"user_id", "chat_id", "content", "created_at"},
		pgx.CopyFromSlice(len(messages), func(i int) ([]any, error) {
			msg := messages[i]
			return []any{msg.UserID, msg.ChatID, msg.Content, msg.Timestamp}, nil
		}),
	)

	if err != nil {
		return fmt.Errorf("%s: copy from failed: %w", op, err)
	}

	return tx.Commit(ctx)
}
