package ports

import "awesome-chat/internal/domain/core/message/entity"

type EntityCreator interface {
	Do(
		userID string,
		chatID string,
		text string,
	) entity.OldMessage
}
