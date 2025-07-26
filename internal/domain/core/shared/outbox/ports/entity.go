package ports

import (
	"awesome-chat/internal/domain/core/shared/outbox/entity"
	"awesome-chat/internal/domain/core/shared/outbox/vo"

	"github.com/google/uuid"
)

type EntityCreator interface {
	Do(
		outboxID uuid.UUID,
		entityName vo.EntityName,
		payload []byte,
	) entity.Outbox
}
