package user

import (
	"awesome-chat/internal/domain/core/shared/ports"
	"context"

	"github.com/google/uuid"
)

type CreateStore struct {
	e ports.ExecutorManager
}

func NewCreateStore(e ports.ExecutorManager) *CreateStore {
	return &CreateStore{e: e}
}

func (s *CreateStore) Execute(ctx context.Context, id uuid.UUID, username string) error {
	conn := s.e.GetPoolExecutor()
	query := `INSERT INTO users (id, username) VALUES ($1, $2)`

	_, err := conn.Exec(ctx, query, id, username)
	if err != nil {
		return err
	}

	return nil
}
