package query

import (
	"context"

	"github.com/google/uuid"

	"awesome-chat/internal/application/user/port"
	"awesome-chat/internal/domain/core/user/entity"
)

type Read struct {
	reader port.Reader
}

func NewRead(reader port.Reader) *Read {
	return &Read{reader: reader}
}

type ReadIn struct {
	UserID uuid.UUID
}

type ReadOut struct {
	User entity.User
}

func (q *Read) Handle(ctx context.Context, in ReadIn) (ReadOut, error) {
	user, err := q.reader.ByID(ctx, in.UserID)
	if err != nil {
		return ReadOut{}, err
	}

	return ReadOut{User: user}, nil
}
