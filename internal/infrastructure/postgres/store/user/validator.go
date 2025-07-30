package user

import (
	sharedPorts "awesome-chat/internal/domain/core/shared/ports"
	userErrors "awesome-chat/internal/domain/core/user/errors"
	"awesome-chat/internal/domain/core/user/ports"
	"awesome-chat/internal/domain/core/user/vo"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

var _ ports.UserValidatorStore = (*ValidatorStore)(nil)

type ValidatorStore struct {
	executor sharedPorts.ExecutorManager
}

func NewValidatorStore(executor sharedPorts.ExecutorManager) *ValidatorStore {
	return &ValidatorStore{executor: executor}
}

func (s *ValidatorStore) ValidateByID(ctx context.Context, userID vo.UserID) error {
	const op = "user.ValidatorStore.ValidateByID"

	id := userID.ToUUID()
	query := `
	SELECT 1 FROM users WHERE id = $1;
	`

	var ok int
	err := s.executor.GetPoolExecutor().QueryRow(ctx, query, id).Scan(&ok)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf("%w: %s", userErrors.ErrUserDoesNotExist, userID)
	case err != nil:
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *ValidatorStore) ValidateByEmail(ctx context.Context, email vo.Email) error {
	const op = "user.ValidatorStore.ValidateByEmail"

	emailStr := email.String()
	query := `
    SELECT 1 FROM users WHERE email = $1;
    `

	var ok int
	err := s.executor.GetPoolExecutor().QueryRow(ctx, query, emailStr).Scan(&ok)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf("%w: %s", userErrors.ErrUserDoesNotExist, emailStr)
	case err != nil:
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *ValidatorStore) ValidateMultiple(ctx context.Context, userIDs vo.UserIDs) error {
	const op = "user.ValidatorStore.ValidateMultiple"

	ids := userIDs.ToUUIDs()
	query := `
    SELECT COUNT(*) FROM users WHERE id = ANY($1);
    `

	var count int
	err := s.executor.GetPoolExecutor().QueryRow(ctx, query, ids).Scan(&count)
	switch {
	case err != nil:
		return fmt.Errorf("%s: %w", op, err)
	case count != len(ids):
		return fmt.Errorf("%s: %w: expected %d, found %d",
			op, userErrors.ErrNotAllUsersExist, len(ids), count,
		)
	}

	return nil
}
