package ports

import (
	"awesome-chat/internal/domain/core/user/entity"
	"awesome-chat/internal/domain/core/user/vo"
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

type UserValidatorStore interface {
	ValidateByID(ctx context.Context, userID vo.UserID) error
	ValidateByEmail(ctx context.Context, email vo.Email) error
	ValidateMultiple(ctx context.Context, userIDs vo.UserIDs) error
}

// var _ Validator = (*ValidatorStore)(nil)
