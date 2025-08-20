package usecases

import (
	"awesome-chat/internal/application/message/dto"
	"context"
)

type SendVoice interface {
	Execute(ctx context.Context, req *dto.SendVoiceRequest) (dto.SendVoiceResponse, error)
}
