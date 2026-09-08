package message

import (
	"context"
	"fmt"

	"awesome-chat/internal/domain/core/message/entity"
	"awesome-chat/internal/domain/core/message/vo"
	"awesome-chat/internal/domain/core/shared/ports"
)

type GetForChatWithFilter struct {
	executor ports.ExecutorManager
}

func NewGetForChatWithFilter(executor ports.ExecutorManager) *GetForChatWithFilter {
	return &GetForChatWithFilter{executor: executor}
}

func (s *GetForChatWithFilter) List(
	ctx context.Context,
	filter vo.ReadFilter,
) ([]entity.MessageForPreview, error) {
	const op = "message.GetForChatWithFilter.List"

	query := `
        SELECT
            id, user_id, content, created_at
        FROM messages
        WHERE chat_id = $1 AND ($2 = 0 OR id < $2)
        ORDER BY id DESC
        LIMIT $3`

	rows, err := s.executor.GetPoolExecutor().Query(ctx, query, filter.ChatID, filter.Cursor, filter.Limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	messages := make([]entity.MessageForPreview, 0, filter.Limit)
	for rows.Next() {
		var message entity.MessageForPreview
		if err = rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.Text,
			&message.Timestamp,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return messages, nil
}
