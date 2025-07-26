package chat

import (
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ValidatorStore struct {
	executor ports.ExecutorManager
}

func NewValidatorStore(executor ports.ExecutorManager) *ValidatorStore {
	return &ValidatorStore{executor: executor}
}

func (v *ValidatorStore) ValidateExists(ctx context.Context, id uuid.UUID) error {
	conn := v.executor.GetPoolExecutor()
	query := `SELECT 1 FROM chat WHERE id = $1 LIMIT 1;`

	var exists int
	if err := conn.QueryRow(ctx, query, id).Scan(&exists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("chat %s not found", id) // 404
		}
		return fmt.Errorf("db error: %w", err) // 500
	}

	return nil
}

func (v *ValidatorStore) IsMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error) {
	conn := v.executor.GetPoolExecutor()
	query := `SELECT 1 FROM user_chats WHERE chat_id = $1 AND user_id = $2 LIMIT 1;`

	var exists int
	if err := conn.QueryRow(ctx, query, chatID, userID).Scan(&exists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("db error: %w", err) // 500
	}

	return true, nil
}
