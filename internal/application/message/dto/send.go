package dto

type (
	Message struct {
		UserID  string `json:"user_id"`
		ChatID  string `json:"chat_id"`
		Content string `json:"content"`
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
}
