package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"awesome-chat/internal/domain/core/shared/ports"
	"awesome-chat/internal/domain/core/user/entity"
	"awesome-chat/internal/domain/core/user/vo"
)

type GetStore struct {
	executor ports.ExecutorManager
}

func NewGetStore(e ports.ExecutorManager) *GetStore {
	return &GetStore{executor: e}
}

func (s *GetStore) All(ctx context.Context) ([]entity.User, error) {
	const op = "user.GetStore.All"

	query := `
	SELECT
	    id, email, password, username, created_at, updated_at
	FROM users
	ORDER BY username`

	rows, err := s.executor.GetPoolExecutor().Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		user, scanErr := scanUser(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("%s: %w", op, scanErr)
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (s *GetStore) ByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	const op = "user.GetStore.ByID"

	query := `
	SELECT
	    id, email, password, username, created_at, updated_at
	FROM users
	WHERE id = $1`

	user, err := scanUser(s.executor.GetPoolExecutor().QueryRow(ctx, query, userID))
	if err != nil {
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(row scanner) (entity.User, error) {
	var (
		user  entity.User
		email string
	)

	if err := row.Scan(
		&user.UserID,
		&email,
		&user.Password,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return entity.User{}, err
	}

	user.Email = vo.Email(email)

	return user, nil
}
