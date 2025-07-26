package user

import (
	"awesome-chat/internal/domain/core/shared/ports"
	"awesome-chat/internal/domain/core/user/entity"
	userErrors "awesome-chat/internal/domain/core/user/errors"
	"awesome-chat/internal/domain/core/user/vo"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type ProviderStore struct {
	executor ports.ExecutorManager
}

func NewProviderStore(executor ports.ExecutorManager) *ProviderStore {
	return &ProviderStore{executor: executor}
}

func (s *ProviderStore) Get(ctx context.Context, email string) (entity.User, error) {
	const op = "postgres.UserProviderStore.Get"

	query := `SELECT 
		id,
		email,
		password,
		username,
		created_at,
		updated_at       			
	FROM users 
	WHERE email = $1`

	conn := s.executor.GetPoolExecutor()

	row := conn.QueryRow(ctx, query, email)

	var user entity.User
	var userEmail string
	err := row.Scan(
		&user.UserID,
		&userEmail,
		&user.Password,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	user.Email = vo.Email(userEmail)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return entity.User{}, fmt.Errorf("%s: %w", op, userErrors.ErrUserDoesNotExist)
	case err != nil:
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	default:
		return user, nil
	}
}
