package getFunc

import (
	"awesome-chat/internal/application/message/dto"
	"awesome-chat/internal/domain/core/message/ports/usecases"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func NewGetMessagesFunc(
	poolGetter func() *pgxpool.Pool,
) usecases.GetMessagesFunc {
	return func(ctx context.Context, req dto.GetRequest) (dto.Messages, error) {
		pool := poolGetter()
		query := `
		SELECT 
			user_id, chat_id, content, created_at
		FROM message 
		WHERE chat_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
		`

		rows, err := pool.Query(ctx, query, req.ChatID, req.Limit, req.Offset)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		messages := make([]dto.Message, 0)
		for rows.Next() {
			var msg dto.Message
			var createdAt time.Time
			if err = rows.Scan(
				&msg.UserID,
				&msg.ChatID,
				&msg.Content,
				&createdAt,
			); err != nil {
				return nil, err
			}
			msg.Timestamp = createdAt.Format(time.RFC3339)
			messages = append(messages, msg)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		return messages, nil
	}
}
