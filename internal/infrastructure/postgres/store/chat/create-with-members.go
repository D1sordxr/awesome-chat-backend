package chat

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"awesome-chat/internal/domain/core/shared/ports"
)

type CreateWithMembersStore struct {
	executor ports.ExecutorManager
}

func NewCreateWithMembersStore(e ports.ExecutorManager) *CreateWithMembersStore {
	return &CreateWithMembersStore{executor: e}
}

func (s *CreateWithMembersStore) Create(
	ctx context.Context,
	chatID uuid.UUID,
	chatName string,
) error {
	const op = "CreateWithMembersStore.Create"

	query := "INSERT INTO chats (id, chat_name) VALUES ($1, $2)"

	conn, err := s.executor.GetTxExecutor(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = conn.Exec(ctx, query, chatID, chatName); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *CreateWithMembersStore) AddMembers(
	ctx context.Context,
	chatID uuid.UUID,
	memberIDs []uuid.UUID,
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

	if _, err = conn.Exec(ctx, query, chatID, memberIDs); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *CreateWithMembersStore) AddMember(ctx context.Context, chatID uuid.UUID, memberID uuid.UUID) error {
	const op = "CreateWithMembersStore.AddMember"

	query := "INSERT INTO user_chats (user_id, chat_id) VALUES ($1, $2)"

	if _, err := s.executor.GetExecutor(ctx).Exec(ctx, query, memberID, chatID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
