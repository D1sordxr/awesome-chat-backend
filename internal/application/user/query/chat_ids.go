package query

import (
	"context"

	"github.com/google/uuid"

	"awesome-chat/internal/application/user/port"
)

type ChatIDs struct {
	reader port.ChatIDsReader
}

func NewChatIDs(reader port.ChatIDsReader) *ChatIDs {
	return &ChatIDs{reader: reader}
}

type ChatIDsIn struct {
	UserID uuid.UUID
}

type ChatIDsOut struct {
	ChatIDs []uuid.UUID
}

func (q *ChatIDs) Handle(ctx context.Context, in ChatIDsIn) (ChatIDsOut, error) {
	chatIDs, err := q.reader.ChatIDs(ctx, in.UserID)
	if err != nil {
		return ChatIDsOut{}, err
	}

	return ChatIDsOut{ChatIDs: chatIDs}, nil
}
