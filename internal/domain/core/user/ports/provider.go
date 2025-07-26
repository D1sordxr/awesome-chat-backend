package ports

import (
	"awesome-chat/internal/domain/core/user/entity"
	"context"
)

type UserProviderStore interface {
	Get(ctx context.Context, email string) (entity.User, error)
}
