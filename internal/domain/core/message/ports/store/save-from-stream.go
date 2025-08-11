package store

import (
	"awesome-chat/internal/domain/core/message/vo"
	"context"
)

type SaveFromStreamStore interface {
	SaveBatch(ctx context.Context, messages []vo.StreamMessage) error
}
