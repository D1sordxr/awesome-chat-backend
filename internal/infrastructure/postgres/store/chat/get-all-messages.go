package chat

import (
	"awesome-chat/internal/domain/core/chat/entity"
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type GetAllMessagesStore struct {
	executor ports.ExecutorManager
}

func NewGetAllMessagesStore(executor ports.ExecutorManager) *GetAllMessagesStore {
	return &GetAllMessagesStore{
		executor: executor,
	}
}

func (s *GetAllMessagesStore) Execute(
	ctx context.Context,
	id uuid.UUID,
) (
	[]entity.MessagePreview,
	error,
) {
	const op = "chat.GetAllMessagesStore.SetupChatPreviews"

	conn := s.executor.GetPoolExecutor()
	query := `
	SELECT user_id, content, created_at 
	FROM messages 
	WHERE chat_id = $1
	ORDER BY created_at DESC
	`

	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var messages []entity.MessagePreview
	for rows.Next() {
		var msg entity.MessagePreview
		if err = rows.Scan(
			&msg.SenderID,
			&msg.Text,
			&msg.Timestamp,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return messages, nil
}
