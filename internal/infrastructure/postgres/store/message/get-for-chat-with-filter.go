package message

import (
	"awesome-chat/internal/domain/core/message/entity"
	"awesome-chat/internal/domain/core/message/vo"
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
	"fmt"
	"strconv"
)

type GetForChatWithFilter struct {
	executor ports.ExecutorManager
}

func NewGetForChatWithFilter(executor ports.ExecutorManager) *GetForChatWithFilter {
	return &GetForChatWithFilter{executor: executor}
}

func (s *GetForChatWithFilter) Execute(
	ctx context.Context,
	filter vo.ReadFilter,
) ([]entity.MessageForPreview, error) {
	const op = "message.GetStore.GetByCursor"

	if filter.Limit <= 0 {
		filter.Limit = 100
	}

	baseQuery := `
        SELECT 
            id, user_id, content, created_at
        FROM messages 
        WHERE chat_id = $1
    `

	var args []interface{}
	args = append(args, filter.ChatID)

	cursorQuery := ""
	if filter.Cursor > 0 {
		cursorQuery = " AND id < $2"
		args = append(args, filter.Cursor)
	}

	fullQuery := baseQuery + cursorQuery + `
        ORDER BY created_at DESC
        LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, filter.Limit)

	rows, err := s.executor.GetPoolExecutor().Query(ctx, fullQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	messages := make([]entity.MessageForPreview, 0, filter.Limit)
	for rows.Next() {
		var msg entity.MessageForPreview
		if err = rows.Scan(
			&msg.ID,
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
