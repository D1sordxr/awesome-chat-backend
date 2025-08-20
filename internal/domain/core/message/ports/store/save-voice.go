package store

import (
	"awesome-chat/internal/domain/core/message/vo"
	"context"
)

type SaveVoice interface {
	Execute(ctx context.Context, data vo.SaveVoiceData) error
}
