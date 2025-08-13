package usecases

import (
	"awesome-chat/internal/application/message/dto"
	"context"
)

type GetMessagesFunc func(ctx context.Context, req dto.GetRequest) (dto.Messages, error)
