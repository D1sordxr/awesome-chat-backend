package message

import (
	"context"
	"fmt"

	"awesome-chat/internal/domain/core/message/entity"
	"awesome-chat/internal/domain/core/message/vo"
	"awesome-chat/internal/domain/core/shared/ports"
)

type SaveVoiceStore struct {
	executor ports.ExecutorManager
}

func NewSaveVoiceStore(executor ports.ExecutorManager) *SaveVoiceStore {
	return &SaveVoiceStore{executor: executor}
}

func (s *SaveVoiceStore) SaveVoice(ctx context.Context, data vo.SaveVoiceData) (entity.Message, error) {
	const op = "message.SaveVoiceStore.SaveVoice"

	tx, err := s.executor.GetTxExecutor(ctx)
	if err != nil {
		return entity.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	messageQuery := `
		INSERT INTO messages (
			user_id,
			chat_id,
			message_type
		) VALUES ($1, $2, 'voice')
		RETURNING id, created_at`

	saved := entity.Message{UserID: data.UserID, ChatID: data.ChatID}

	if err = tx.QueryRow(ctx, messageQuery, data.UserID, data.ChatID).Scan(
		&saved.ID,
		&saved.Timestamp,
	); err != nil {
		return entity.Message{}, fmt.Errorf("%s: insert message: %w", op, err)
	}

	voiceQuery := `
		INSERT INTO voice_messages (
			message_id,
			object_key,
			duration_seconds,
			waveform
		) VALUES ($1, $2, $3, $4)`

	if _, err = tx.Exec(
		ctx,
		voiceQuery,
		saved.ID,
		data.ObjectKey,
		data.DurationSeconds,
		data.Waveform,
	); err != nil {
		return entity.Message{}, fmt.Errorf("%s: insert voice message: %w", op, err)
	}

	return saved, nil
}
