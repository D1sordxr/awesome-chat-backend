package sendVoice

import (
	"awesome-chat/internal/application/message/dto"
	"context"
)

type MessageSendVoiceUseCase struct {
	//
}

func NewMessageSendVoiceUseCase() *MessageSendVoiceUseCase {
	return &MessageSendVoiceUseCase{}
}

func (uc *MessageSendVoiceUseCase) Execute(
	ctx context.Context,
	req dto.SendVoiceRequest,
) (
	dto.SendVoiceResponse,
	error,
) {
	return dto.SendVoiceResponse{}, nil
}
