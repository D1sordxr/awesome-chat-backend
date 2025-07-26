package filters

import (
	"awesome-chat/internal/domain/core/shared/outbox/vo"

	"github.com/google/uuid"
)

type SetOutboxStatus struct {
	OutboxID uuid.UUID
	Status   vo.OutboxStatus
}
