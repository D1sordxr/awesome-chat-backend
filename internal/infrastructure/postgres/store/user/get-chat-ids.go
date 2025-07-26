package user

import (
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
	"github.com/google/uuid"
)

type GetChatIDsStore struct {
	e ports.ExecutorManager
}

func NewGetChatIDsStore(e ports.ExecutorManager) *GetChatIDsStore {
	return &GetChatIDsStore{e: e}
}

func (s *GetChatIDsStore) Execute(ctx context.Context, userID uuid.UUID) ([]string, error) {
	query := `SELECT chat_id FROM user_chats WHERE user_id = $1`

	rows, err := s.e.GetPoolExecutor().Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chatIDs := make([]uuid.UUID, 0, 10)
	for rows.Next() {
		var scanID uuid.UUID
		if err = rows.Scan(&scanID); err != nil {
			return nil, err
		}
		chatIDs = append(chatIDs, scanID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	chatIDsStr := make([]string, 0, len(chatIDs))
	for _, chatID := range chatIDs {
		chatIDsStr = append(chatIDsStr, chatID.String())
	}

	return chatIDsStr, nil
}
