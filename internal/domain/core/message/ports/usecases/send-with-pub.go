package usecases

import (
	"awesome-chat/internal/application/message/dto"
	"context"
)

type MessageBroadcastWithPub interface {
	Execute(ctx context.Context, req dto.BroadcastWithPubRequest) error
}
