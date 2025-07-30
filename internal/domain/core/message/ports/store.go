package ports

import (
	"awesome-chat/internal/domain/core/message/entity"
	"awesome-chat/internal/domain/core/message/vo"
	"context"
)

type GetMessageStore interface {
	GetAllMessagesFromChat(ctx context.Context, chatID string) ([]entity.OldMessage, error)
}
type GetForChatWithFilterStore interface {
	Execute(ctx context.Context, filter vo.ReadFilter) ([]entity.MessageForPreview, error)
}
