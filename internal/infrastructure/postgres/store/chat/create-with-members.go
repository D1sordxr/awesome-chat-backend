package chat

import (
	"awesome-chat/internal/domain/core/shared/ports"
	"context"

	"github.com/google/uuid"
)

type CreateWithMembersStore struct {
	e ports.ExecutorManager
}

func NewCreateWithMembersStore(e ports.ExecutorManager) *CreateWithMembersStore {
	return &CreateWithMembersStore{e: e}
}

func (cs *CreateWithMembersStore) Execute(ctx context.Context, id uuid.UUID, name string) error {
	conn, err := cs.e.GetTxExecutor(ctx)
	if err != nil {
		return err
	}
	query := "INSERT INTO chats (id, chat_name) VALUES ($1, $2)"

	if _, err = conn.Exec(ctx, query, id, name); err != nil {
		return err
	}

	return nil
}

func (cs *CreateWithMembersStore) AddMember(ctx context.Context, chatID uuid.UUID, memberID uuid.UUID) error {
	conn := cs.e.GetExecutor(ctx)
	query := "INSERT INTO user_chats (user_id, chat_id) VALUES ($1, $2)"

	if _, err := conn.Exec(ctx, query, memberID, chatID); err != nil {
		return err
	}

	return nil
}
