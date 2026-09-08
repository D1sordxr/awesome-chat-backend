package port

import (
	"context"

	"github.com/google/uuid"

	"awesome-chat/internal/domain/core/user/entity"
)

//go:generate mockgen -source=port.go -destination=mocks/port_mocks.go -package=mocks

type Repository interface {
	Save(ctx context.Context, user entity.User) error
}

type Provider interface {
	ByEmail(ctx context.Context, email string) (entity.User, error)
}

type Reader interface {
	ByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
	All(ctx context.Context) ([]entity.User, error)
}

type ChatIDsReader interface {
	ChatIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

type TokenCreator interface {
	Do(user entity.User) (string, error)
}
