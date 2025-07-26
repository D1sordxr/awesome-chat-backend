package ports

import (
	"awesome-chat/internal/domain/core/message/entity"
	"context"
)

type Repository interface {
	Save(ctx context.Context, message entity.Message) error
}
