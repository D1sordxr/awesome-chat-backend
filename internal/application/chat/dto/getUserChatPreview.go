package dto

type (
	UserID string

	GetUserChatPreviewResponse struct {
		ChatPreviews []ChatPreview `json:"chat_previews"`
	}
	ChatPreview struct {
		ChatID       string        `json:"chat_id"`
		Name         string        `json:"name"`
		LastMessage  Message       `json:"last_message,omitempty"`
		UnreadCount  int           `json:"unread_count"`
		AvatarURL    string        `json:"avatar_url,omitempty"`
		Participants []Participant `json:"participants,omitempty"`
	}
	Message struct {
		UserID    string `json:"user_id"`
		Content   string `json:"content"`
		Timestamp string `json:"timestamp"`
	}
	Participant struct {
		UserID    string `json:"user_id"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url,omitempty"`
		IsOnline  bool   `json:"is_online,omitempty"`
	}
)
