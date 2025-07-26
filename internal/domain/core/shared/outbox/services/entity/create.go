package entity

import (
	"awesome-chat/internal/domain/core/shared/outbox/entity"
	"awesome-chat/internal/domain/core/shared/outbox/vo"

	"github.com/google/uuid"
)

type Create struct{}

func (*Create) Do(
	id uuid.UUID,
	entityName vo.EntityName,
	payload []byte,
) entity.Outbox {
	return entity.Outbox{
		OutboxID:   id,
		EntityName: entityName,
		Payload:    payload,
		Status:     vo.StatusPending,
	}
}
