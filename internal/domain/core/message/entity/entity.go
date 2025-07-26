package entity

type Message struct {
	UserID  string `json:"user_id"`
	ChatID  string `json:"chat_id"`
	Content string `json:"data"`
}
