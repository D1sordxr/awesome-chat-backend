package ports

import (
	"awesome-chat/internal/domain/core/chat/entity"
	"awesome-chat/internal/domain/core/chat/vo"
	userVO "awesome-chat/internal/domain/core/user/vo"
	"context"

	"github.com/google/uuid"
)

type CreateWithMembersStore interface {
	CreateChat(ctx context.Context, chatID vo.ChatID, chatName string) error
	AddMembers(ctx context.Context, chatID vo.ChatID, memberIDs userVO.UserIDs) error
}

type GetAllMessagesStore interface {
	Execute(ctx context.Context, id uuid.UUID) ([]entity.MessagePreview, error)
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
