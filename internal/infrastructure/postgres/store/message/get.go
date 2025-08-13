package message

import (
	"awesome-chat/internal/domain/core/message/entity"
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
)

type GetStore struct {
	executor ports.ExecutorManager
}

func NewGetStore(executor ports.ExecutorManager) *GetStore {
	return &GetStore{executor: executor}
}

func (s *GetStore) GetAllMessagesFromChat(ctx context.Context, chatID string) ([]entity.OldMessage, error) {
	conn := s.executor.GetPoolExecutor()
	query := `SELECT 
		user_id, chat_id, content
	FROM message WHERE chat_id = $1
	ORDER BY created_at DESC`

	rows, err := conn.Query(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]entity.OldMessage, 0, 20)
	for rows.Next() {
		var m entity.OldMessage
		if err = rows.Scan(
			&m.UserID,
			&m.ChatID,
			&m.Content,
		); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}
