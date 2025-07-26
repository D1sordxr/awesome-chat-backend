package ports

import (
	"awesome-chat/internal/domain/core/user/entity"
	"context"

	"github.com/google/uuid"
)

type UserCreateStore interface {
	Execute(ctx context.Context, id uuid.UUID, username string) error
}

type UserGetAllStore interface {
	GetAll(ctx context.Context) ([]entity.UserRaw, error)
}

type UserGetStore interface {
	Execute(ctx context.Context, id uuid.UUID) (entity.User, error)
}
