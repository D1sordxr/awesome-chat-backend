package fast

import (
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
)

type MessageSendUseCase struct {
	pub ports.Publisher
}

func NewMessageSendFastUseCase(pub ports.Publisher) *MessageSendUseCase {
	return &MessageSendUseCase{
		pub: pub,
	}
}

func (uc *MessageSendUseCase) Execute(ctx context.Context, payload []byte) error {
	return uc.pub.Publish(ctx, payload)
}
