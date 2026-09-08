package vo

import "github.com/google/uuid"

type ReadFilter struct {
	ChatID uuid.UUID `json:"chat_id"`
	Limit  int       `json:"limit,omitempty"`
	Cursor int64     `json:"cursor,omitempty"`
}
