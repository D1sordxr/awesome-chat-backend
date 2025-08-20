package vo

import "github.com/google/uuid"

type SaveVoiceData struct {
	MessageID int64
	UserID    uuid.UUID
	ChatID    uuid.UUID
	AudioURL  string
	Duration  int
	Waveform  []byte
}
