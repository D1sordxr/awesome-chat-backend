package query

import (
	"context"

	"github.com/google/uuid"

	"awesome-chat/internal/application/message/port"
	chatErrors "awesome-chat/internal/domain/core/chat/errors"
	"awesome-chat/internal/domain/core/message/entity"
	"awesome-chat/internal/domain/core/message/vo"
)

const (
	defaultLimit = 100
	maxLimit     = 500
)

type List struct {
	reader port.Reader
	chats  port.ChatValidator
}

func NewList(reader port.Reader, chats port.ChatValidator) *List {
	return &List{reader: reader, chats: chats}
}

type ListIn struct {
	CallerID uuid.UUID
	ChatID   uuid.UUID
	Limit    int
	Cursor   int64
}

type ListOut struct {
	Messages   []entity.MessageForPreview
	NextCursor int64
	HasMore    bool
}

func (q *List) Handle(ctx context.Context, in ListIn) (ListOut, error) {
	isMember, err := q.chats.IsMember(ctx, in.ChatID, in.CallerID)
	if err != nil {
		return ListOut{}, err
	}

	if !isMember {
		return ListOut{}, chatErrors.ErrChatAccessDenied
	}

	limit := in.Limit
	switch {
	case limit <= 0:
		limit = defaultLimit
	case limit > maxLimit:
		limit = maxLimit
	}

	messages, err := q.reader.List(ctx, vo.ReadFilter{
		ChatID: in.ChatID,
		Limit:  limit,
		Cursor: in.Cursor,
	})
	if err != nil {
		return ListOut{}, err
	}

	out := ListOut{Messages: messages, HasMore: len(messages) == limit}
	if len(messages) > 0 {
		out.NextCursor = messages[len(messages)-1].ID
	}

	return out, nil
}
