package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"awesome-chat/internal/domain/core/shared/ports"
	"awesome-chat/internal/domain/core/user/entity"
	userErrors "awesome-chat/internal/domain/core/user/errors"
)

type ProviderStore struct {
	executor ports.ExecutorManager
}

func NewProviderStore(executor ports.ExecutorManager) *ProviderStore {
	return &ProviderStore{executor: executor}
}

func (s *ProviderStore) ByEmail(ctx context.Context, email string) (entity.User, error) {
	const op = "user.ProviderStore.ByEmail"

	query := `
	SELECT
		id, email, password, username, created_at, updated_at
	FROM users
	WHERE email = $1`

	user, err := scanUser(s.executor.GetPoolExecutor().QueryRow(ctx, query, email))

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return entity.User{}, fmt.Errorf("%s: %w", op, userErrors.ErrUserDoesNotExist)
	case err != nil:
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	default:
		return user, nil
	}
}
