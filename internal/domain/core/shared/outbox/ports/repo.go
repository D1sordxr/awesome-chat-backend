package ports

import (
	"awesome-chat/internal/domain/core/shared/outbox/entity"
	"context"
)

type Repository interface {
	Save(ctx context.Context, outbox entity.Outbox) error
}
