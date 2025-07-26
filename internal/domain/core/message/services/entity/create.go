package entity

import "awesome-chat/internal/domain/core/message/entity"

type Create struct{}

func (*Create) Do(
	userID string,
	chatID string,
	text string,
) entity.Message {
	return entity.Message{
		UserID:  userID,
		ChatID:  chatID,
		Content: text,
	}
}
