package entity

import (
	"github.com/google/uuid"
	"time"
)

type (
	ChatPreview struct {
		ChatID       uuid.UUID      `json:"chat_id"`
		Name         string         `json:"name"`
		LastMessage  MessagePreview `json:"last_message,omitempty"`
		UnreadCount  int            `json:"unread_count,omitempty"`
		AvatarURL    string         `json:"avatar_url,omitempty"`
		Participants []Participant  `json:"participants"`
	}
	MessagePreview struct {
		ID        int       `json:"id"`
		SenderID  uuid.UUID `json:"sender_id"`
		Text      string    `json:"text"`
		Timestamp time.Time `json:"timestamp"`
	}
	Participant struct {
		UserID    uuid.UUID `json:"user_id"`
		Username  string    `json:"username"`
		AvatarURL string    `json:"avatar_url,omitempty"`
		IsOnline  bool      `json:"is_online,omitempty"`
	}
)
