package ports

import (
	"awesome-chat/internal/domain/core/chat/entity"
	"context"

	"github.com/google/uuid"
)

type CreateWithMembersStore interface {
	CreateStore
	AddMemberStore
}

type GetAllMessagesStore interface {
	Execute(ctx context.Context, id uuid.UUID) ([]entity.MessagePreview, error)
}

type CreateStore interface {
	Execute(ctx context.Context, id uuid.UUID, chatName string) error
}

type AddMemberStore interface {
	AddMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
}

type ValidateStore interface {
	ValidateExists(ctx context.Context, id uuid.UUID) error
	IsMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error)
}
type GetUserChatPreviewStore interface {
	SetupChatPreviews(
		ctx context.Context,
		userID uuid.UUID,
	) (
		previews []entity.ChatPreview,
		err error,
	)
	BatchParticipantsForChats(
		ctx context.Context,
		chatIDs []uuid.UUID,
	) (
		map[uuid.UUID][]entity.Participant,
		error,
	)
}
