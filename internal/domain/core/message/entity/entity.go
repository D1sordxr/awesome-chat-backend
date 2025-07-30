package entity

import (
	"github.com/google/uuid"
	"time"
)

type OldMessage struct {
	ID      int    `json:"id"`
	UserID  string `json:"user_id"`
	ChatID  string `json:"chat_id"`
	Content string `json:"data"`
}

type MessageForPreview struct {
	ID        int       `json:"id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}
