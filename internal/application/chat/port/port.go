package port

import (
	"context"

	"github.com/google/uuid"

	"awesome-chat/internal/domain/core/chat/entity"
)

//go:generate mockgen -source=port.go -destination=mocks/port_mocks.go -package=mocks

type Repository interface {
	Create(ctx context.Context, chatID uuid.UUID, name string) error
	AddMembers(ctx context.Context, chatID uuid.UUID, memberIDs []uuid.UUID) error
	AddMember(ctx context.Context, chatID uuid.UUID, memberID uuid.UUID) error
}

type Validator interface {
	Exists(ctx context.Context, chatID uuid.UUID) error
	IsMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error)
}

type PreviewReader interface {
	Previews(ctx context.Context, userID uuid.UUID) ([]entity.ChatPreview, error)
	BatchParticipants(ctx context.Context, chatIDs []uuid.UUID) (map[uuid.UUID][]entity.Participant, error)
}

type UserValidator interface {
	ValidateByID(ctx context.Context, userID uuid.UUID) error
	ValidateMultiple(ctx context.Context, userIDs []uuid.UUID) error
}

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
