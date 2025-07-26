package user

import (
	"awesome-chat/internal/domain/core/shared/ports"
	"awesome-chat/internal/domain/core/user/entity"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type GetStore struct {
	executor ports.ExecutorManager
}

func NewGetStore(e ports.ExecutorManager) *GetStore {
	return &GetStore{executor: e}
}

func (s *GetStore) GetAll(ctx context.Context) ([]entity.UserRaw, error) {
	const op = "user.GetStore.GetAll"
	query := `
	SELECT 
	    id, email, password, username, created_at, updated_at
	FROM users
	`

	rows, err := s.executor.GetPoolExecutor().Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var users []entity.UserRaw
	for rows.Next() {
		var user entity.UserRaw
		if err = rows.Scan(
			&user.UserID,
			&user.Email,
			&user.Password,
			&user.Username,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return users, nil
}

func (s *GetStore) Execute(ctx context.Context, id uuid.UUID) (entity.User, error) {
	conn := s.executor.GetPoolExecutor()
	query := `SELECT username, created_at, updated_at FROM users WHERE id = $1`

	var user entity.User
	if err := conn.QueryRow(ctx, query, id).Scan(
		&user.Username, &user.CreatedAt, &user.UpdatedAt,
	); err != nil {
		return entity.User{}, err
	}

	return user, nil
}
