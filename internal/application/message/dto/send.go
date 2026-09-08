package dto

import "time"

type (
	Message struct {
		UserID    string    `json:"user_id"`
		ChatID    string    `json:"chat_id"`
		Content   string    `json:"content"`
		Timestamp time.Time `json:"timestamp,omitempty"`
	}
	BroadcastWithPubRequest struct {
		Message
	}
)
