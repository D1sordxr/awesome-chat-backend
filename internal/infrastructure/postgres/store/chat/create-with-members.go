package chat

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"awesome-chat/internal/domain/core/chat/vo"
	"awesome-chat/internal/domain/core/shared/ports"
	userVO "awesome-chat/internal/domain/core/user/vo"
)

type CreateWithMembersStore struct {
	executor ports.ExecutorManager
}

func NewCreateWithMembersStore(e ports.ExecutorManager) *CreateWithMembersStore {
	return &CreateWithMembersStore{executor: e}
}

func (s *CreateWithMembersStore) CreateChat(
	ctx context.Context,
	chatID vo.ChatID,
	chatName string,
) error {
	const op = "CreateWithMembersStore.CreateChat"

	query := "INSERT INTO chats (id, chat_name) VALUES ($1, $2)"

	id := chatID.ToUUID()
	conn, err := s.executor.GetTxExecutor(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = conn.Exec(ctx, query, id, chatName); err != nil {
		return err
	}

	return nil
}

func (s *CreateWithMembersStore) AddMembers(
	ctx context.Context,
	chatID vo.ChatID,
	memberIDs userVO.UserIDs,
) error {
	const op = "CreateWithMembersStore.AddMembers"
	query := `
    INSERT INTO user_chats (chat_id, user_id) 
    SELECT $1, unnest($2::uuid[])
    ON CONFLICT DO NOTHING
    `

	conn, err := s.executor.GetTxExecutor(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	chatUUID := chatID.ToUUID()
	memberUUIDs := memberIDs.ToUUIDs()

	if _, err = conn.Exec(ctx, query, chatUUID, memberUUIDs); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *CreateWithMembersStore) AddMember(ctx context.Context, chatID uuid.UUID, memberID uuid.UUID) error {
	conn := s.executor.GetExecutor(ctx)
	query := "INSERT INTO user_chats (user_id, chat_id) VALUES ($1, $2)"

	if _, err := conn.Exec(ctx, query, memberID, chatID); err != nil {
		return err
	}

	return nil
}
