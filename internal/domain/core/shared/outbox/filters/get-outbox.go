package filters

import "awesome-chat/internal/domain/core/shared/outbox/vo"

type GetOutbox struct {
	EntityName vo.EntityName
	Status     vo.OutboxStatus
	Limit      int
}
