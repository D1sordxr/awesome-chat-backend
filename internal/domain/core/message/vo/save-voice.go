package vo

import "github.com/google/uuid"

type SaveVoiceData struct {
	UserID          uuid.UUID
	ChatID          uuid.UUID
	ObjectKey       string
	DurationSeconds int
	Waveform        []byte
}
