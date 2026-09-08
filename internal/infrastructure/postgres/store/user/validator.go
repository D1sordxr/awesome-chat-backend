package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"awesome-chat/internal/domain/core/shared/ports"
	userErrors "awesome-chat/internal/domain/core/user/errors"
)

type ValidatorStore struct {
	executor ports.ExecutorManager
}

func NewValidatorStore(executor ports.ExecutorManager) *ValidatorStore {
	return &ValidatorStore{executor: executor}
}

func (s *ValidatorStore) ValidateByID(ctx context.Context, userID uuid.UUID) error {
	const op = "user.ValidatorStore.ValidateByID"

	query := `SELECT 1 FROM users WHERE id = $1`

	var ok int
	err := s.executor.GetPoolExecutor().QueryRow(ctx, query, userID).Scan(&ok)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf("%w: %s", userErrors.ErrUserDoesNotExist, userID)
	case err != nil:
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *ValidatorStore) ValidateMultiple(ctx context.Context, userIDs []uuid.UUID) error {
	const op = "user.ValidatorStore.ValidateMultiple"

	query := `SELECT COUNT(*) FROM users WHERE id = ANY($1)`

	var count int
	err := s.executor.GetPoolExecutor().QueryRow(ctx, query, userIDs).Scan(&count)

	switch {
	case err != nil:
		return fmt.Errorf("%s: %w", op, err)
	case count != len(userIDs):
		return fmt.Errorf("%s: %w: expected %d, found %d",
			op, userErrors.ErrNotAllUsersExist, len(userIDs), count,
		)
	}

	return nil
}
