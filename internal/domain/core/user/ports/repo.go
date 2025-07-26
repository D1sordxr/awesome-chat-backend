package ports

import (
	"awesome-chat/internal/domain/core/user/entity"
	"context"
)

type UserRepo interface {
	Save(ctx context.Context, user entity.User) error
}
