package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"awesome-chat/internal/domain/core/shared/ports"
)

type GetChatIDsStore struct {
	e ports.ExecutorManager
}

func NewGetChatIDsStore(e ports.ExecutorManager) *GetChatIDsStore {
	return &GetChatIDsStore{e: e}
}

func (s *GetChatIDsStore) ChatIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	const op = "user.GetChatIDsStore.ChatIDs"

	query := `SELECT chat_id FROM user_chats WHERE user_id = $1`

	rows, err := s.e.GetPoolExecutor().Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	chatIDs := make([]uuid.UUID, 0, 10)
	for rows.Next() {
		var chatID uuid.UUID
		if err = rows.Scan(&chatID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		chatIDs = append(chatIDs, chatID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return chatIDs, nil
}
