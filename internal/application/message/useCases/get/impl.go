package get

import (
	"awesome-chat/internal/application/message/dto"
	"awesome-chat/internal/domain/core/message/ports"
	"context"
)

type MessageGetUseCase struct {
	store ports.GetMessageStore
}

func NewMessageGetUseCase(
	store ports.GetMessageStore,
) *MessageGetUseCase {
	return &MessageGetUseCase{
		store: store,
	}
}

func (uc *MessageGetUseCase) Execute(ctx context.Context, req dto.GetRequest) (dto.Messages, error) {
	entities, err := uc.store.GetAllMessagesFromChat(ctx, req.ChatID)
	if err != nil {
		return nil, err
	}

	messages := make([]dto.Message, 0, len(entities))
	for i := 0; i < len(entities); i++ {
		messages[i] = dto.Message{
			UserID:  entities[i].UserID,
			ChatID:  entities[i].ChatID,
			Content: entities[i].Content,
		}
	}

	return messages, nil
}
