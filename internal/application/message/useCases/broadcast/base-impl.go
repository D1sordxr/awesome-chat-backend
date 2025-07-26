package broadcast

import (
	"awesome-chat/internal/application/message/dto"
	"awesome-chat/internal/domain/core/message/ports"
	wsPorts "awesome-chat/internal/domain/core/shared/ports/ws"
	"context"
	"encoding/json"
)

type MessageBroadcastUseCase struct {
	creator ports.EntityCreator
	br      wsPorts.Broadcaster
}

func NewMessageBroadcastUseCase(
	creator ports.EntityCreator,
	br wsPorts.Broadcaster,
) *MessageBroadcastUseCase {
	return &MessageBroadcastUseCase{
		creator: creator,
		br:      br,
	}
}

func (uc *MessageBroadcastUseCase) Execute(ctx context.Context, message dto.Message) error {
	entity := uc.creator.Do(
		message.UserID,
		message.ChatID,
		message.Content,
	)
	payload, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	if err = uc.br.Broadcast(
		ctx,
		entity,
		payload,
	); err != nil {
		return err
	}

	return nil
}
