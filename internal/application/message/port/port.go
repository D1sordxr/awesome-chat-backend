package port

import (
	"context"
	"time"

	"github.com/google/uuid"

	"awesome-chat/internal/domain/core/message/entity"
	"awesome-chat/internal/domain/core/message/vo"
	outboxEntity "awesome-chat/internal/domain/core/shared/outbox/entity"
)

//go:generate mockgen -source=port.go -destination=mocks/port_mocks.go -package=mocks

type Repository interface {
	Save(ctx context.Context, message entity.Message) (entity.Message, error)
}

type VoiceRepository interface {
	SaveVoice(ctx context.Context, data vo.SaveVoiceData) (entity.Message, error)
}

type Outbox interface {
	Save(ctx context.Context, outbox outboxEntity.Outbox) error
}

type Reader interface {
	List(ctx context.Context, filter vo.ReadFilter) ([]entity.MessageForPreview, error)
}

type UploadService interface {
	PresignedPut(ctx context.Context, objectKey string) (string, time.Time, error)
}

type ChatValidator interface {
	IsMember(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error)
}

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
