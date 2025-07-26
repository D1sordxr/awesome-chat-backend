package chathub

import "encoding/json"

type Message struct {
	UserID   string `json:"user_id,omitempty"`
	ChatID   string `json:"chat_id,omitempty"`
	Content  string `json:"content,omitempty"`
	ServerIP string `json:"server_ip,omitempty"` // k8s
	SenderIP string `json:"sender_ip,omitempty"` // k8s
}

func (m *Message) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}
