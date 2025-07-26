package broadcast

import (
	"awesome-chat/internal/domain/core/message/entity"
	brokerMsg "awesome-chat/internal/domain/core/shared/broker/entity"
	wsPorts "awesome-chat/internal/domain/core/shared/ports/ws"
	"context"
	"encoding/json"
)

type UseCase struct {
	br wsPorts.Broadcaster
}

func NewUseCase(
	br wsPorts.Broadcaster,
) *UseCase {
	return &UseCase{
		br: br,
	}
}

func (uc *UseCase) Broadcast(ctx context.Context, message brokerMsg.Message) error {
	// TODO: save to processed msg by message.key with tx if needed

	var msg entity.Message
	err := json.Unmarshal(message.Value, &msg)
	if err != nil {
		return err
	}
	if err = uc.br.Broadcast(
		ctx,
		msg,
		message.Value,
	); err != nil {
		return err
	}

	return nil
}
