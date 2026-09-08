package query

import (
	"context"

	"awesome-chat/internal/application/user/port"
	"awesome-chat/internal/domain/core/user/entity"
)

type List struct {
	reader port.Reader
}

func NewList(reader port.Reader) *List {
	return &List{reader: reader}
}

type ListIn struct{}

type ListOut struct {
	Users []entity.User
}

func (q *List) Handle(ctx context.Context, _ ListIn) (ListOut, error) {
	users, err := q.reader.All(ctx)
	if err != nil {
		return ListOut{}, err
	}

	return ListOut{Users: users}, nil
}
