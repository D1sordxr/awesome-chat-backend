package chathub

import "encoding/json"

type Message struct {
	UserID    string `json:"user_id"`
	ChatID    string `json:"chat_id"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp,omitempty"`
	ServerIP  string `json:"server_ip,omitempty"` // k8s
	SenderIP  string `json:"sender_ip,omitempty"` // k8s
}

func (m *Message) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}
