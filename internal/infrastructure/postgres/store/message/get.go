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

func (uc *GetStore) GetAllMessagesFromChat(ctx context.Context, chatID string) ([]entity.Message, error) {
	conn := uc.executor.GetPoolExecutor()
	query := `SELECT 
		user_id, chat_id, content
	FROM message WHERE chat_id = $1`

	rows, err := conn.Query(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]entity.Message, 0, 20)
	for rows.Next() {
		var m entity.Message
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
