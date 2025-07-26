package entity

import (
	"awesome-chat/internal/domain/core/shared/outbox/vo"

	"github.com/google/uuid"
)

type Outbox struct {
	OutboxID   uuid.UUID
	EntityName vo.EntityName
	Status     vo.OutboxStatus
	Payload    []byte
}
