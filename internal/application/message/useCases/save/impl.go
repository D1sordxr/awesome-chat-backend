package save

import (
	"awesome-chat/internal/application/message/dto"
	"awesome-chat/internal/domain/core/message/ports"
	"context"
)

type MessageSaveUseCase struct {
	entityCreator ports.EntityCreator
	msgRepo       ports.Repository
}

func NewMessageSaveUseCase(
	entityCreator ports.EntityCreator,
	msgRepo ports.Repository,
) *MessageSaveUseCase {
	return &MessageSaveUseCase{
		entityCreator: entityCreator,
		msgRepo:       msgRepo,
	}
}

func (uc *MessageSaveUseCase) Execute(ctx context.Context, msg dto.Message) error {
	entity := uc.entityCreator.Do(
		msg.UserID,
		msg.ChatID,
		msg.Content,
	)

	return uc.msgRepo.Save(ctx, entity)
}
