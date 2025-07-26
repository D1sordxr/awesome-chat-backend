package ports

import (
	"awesome-chat/internal/domain/core/shared/outbox/entity"
	"awesome-chat/internal/domain/core/shared/outbox/filters"
	"context"
)

type ProcessStore interface {
	GetMessagesForUpdate(ctx context.Context, outbox filters.GetOutbox) ([]entity.Outbox, error)
	UpdateStatus(ctx context.Context, status filters.SetOutboxStatus) error
}
