package message

import (
	"awesome-chat/internal/domain/core/message/vo"
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
	"fmt"
)

type SaveVoiceStore struct {
	executor ports.ExecutorManager
}

func NewSaveVoiceStore(executor ports.ExecutorManager) *SaveVoiceStore {
	return &SaveVoiceStore{executor: executor}
}

func (s *SaveVoiceStore) Execute(ctx context.Context, data vo.SaveVoiceData) error {
	const op = "message.SaveVoiceStore.Execute"

	tx, err := s.executor.GetTxExecutor(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	messageQuery := `
		INSERT INTO messages (
			user_id, 
			chat_id, 
			message_type,
			content
		) VALUES ($1, $2, 'voice', $3)
		RETURNING id
	`

	var messageID int64
	err = tx.QueryRow(ctx, messageQuery, data.UserID, data.ChatID, data.AudioURL).Scan(&messageID)
	if err != nil {
		return fmt.Errorf("%s: failed to insert message: %w", op, err)
	}

	voiceQuery := `
		INSERT INTO voice_messages (
			message_id,
			audio_url,
			duration,
			waveform
		) VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(ctx, voiceQuery, messageID, data.AudioURL, data.Duration, data.Waveform)
	if err != nil {
		return fmt.Errorf("%s: failed to insert voice message: %w", op, err)
	}

	return nil
}
