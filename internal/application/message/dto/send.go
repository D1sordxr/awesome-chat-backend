package dto

type (
	Message struct {
		UserID    string `json:"user_id"`
		ChatID    string `json:"chat_id"`
		Content   string `json:"content"`
		Timestamp string `json:"timestamp,omitempty"`
	}
	Messages []Message
)
type SendRequest struct {
	UserID  string `json:"user_id"`
	ChatID  string `json:"chat_id"`
	Content string `json:"content"`
}

type SendSyncRequest struct {
	UserID  string `json:"user_id"`
	ChatID  string `json:"chat_id"`
	Content string `json:"content"`
}

type GetRequest struct {
	ChatID string `json:"chat_id"`
	Limit  int    `json:"limit,omitempty"`  // 100
	Offset int    `json:"offset,omitempty"` // pagination
	// or
	Cursor string `json:"cursor,omitempty"` // last message ID (cursor pagination)
}
