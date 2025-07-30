package vo

import "github.com/google/uuid"

type ReadFilter struct {
	ChatID uuid.UUID `json:"chat_id"`
	Limit  int       `json:"limit,omitempty"`
	Offset int       `json:"offset,omitempty"`
	Cursor int       `json:"cursor,omitempty"`
}
