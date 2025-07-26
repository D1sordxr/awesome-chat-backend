package repositories

import (
	"awesome-chat/internal/domain/core/shared/ports"
	"awesome-chat/internal/domain/core/user/entity"
	"context"
	"fmt"
)

type UserRepo struct {
	e ports.ExecutorManager
}

func NewUserRepo(e ports.ExecutorManager) *UserRepo {
	return &UserRepo{e: e}
}

func (r *UserRepo) Save(ctx context.Context, user entity.User) error {
	const op = "postgres.UserRepo.Save"

	conn := r.e.GetExecutor(ctx)
	query := `INSERT INTO users (
                  id,
                  email,
                  password,
                  username,
                  created_at,
                  updated_at
			) VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := conn.Exec(ctx, query,
		user.UserID,
		user.Email.String(),
		user.Password,
		user.Username,
		user.CreatedAt,
		user.UpdatedAt,
	); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
