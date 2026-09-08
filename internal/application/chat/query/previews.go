package query

import (
	"context"

	"github.com/google/uuid"

	"awesome-chat/internal/application/chat/port"
	"awesome-chat/internal/domain/core/chat/entity"
)

type Previews struct {
	reader port.PreviewReader
}

func NewPreviews(reader port.PreviewReader) *Previews {
	return &Previews{reader: reader}
}

type PreviewsIn struct {
	UserID uuid.UUID
}

type PreviewsOut struct {
	Chats []entity.ChatPreview
}

func (q *Previews) Handle(ctx context.Context, in PreviewsIn) (PreviewsOut, error) {
	previews, err := q.reader.Previews(ctx, in.UserID)
	if err != nil {
		return PreviewsOut{}, err
	}

	if len(previews) == 0 {
		return PreviewsOut{}, nil
	}

	chatIDs := make([]uuid.UUID, 0, len(previews))
	for _, preview := range previews {
		chatIDs = append(chatIDs, preview.ChatID)
	}

	participants, err := q.reader.BatchParticipants(ctx, chatIDs)
	if err != nil {
		return PreviewsOut{}, err
	}

	for i := range previews {
		previews[i].Participants = participants[previews[i].ChatID]
	}

	return PreviewsOut{Chats: previews}, nil
}
