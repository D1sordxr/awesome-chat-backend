package ports

import (
	"awesome-chat/internal/domain/core/message/entity"
	"context"
)

type GetMessageStore interface {
	GetAllMessagesFromChat(ctx context.Context, chatID string) ([]entity.Message, error)
}
