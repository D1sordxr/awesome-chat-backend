package vo

import (
	"fmt"
	"time"
)

const SentMessageEvent = "sent-message-event"

type (
	StreamMessage struct {
		AckID     string    `json:"ack_id"`
		Event     string    `json:"event"`
		UserID    string    `json:"user_id"`
		ChatID    string    `json:"chat_id"`
		Content   string    `json:"content"`
		Timestamp time.Time `json:"timestamp"`
	}
)

const (
	eventMapKey   = "event"
	userIDMapKey  = "user_id"
	chatIDMapKey  = "chat_id"
	contentMapKey = "content"
	timestampKey  = "timestamp"
)

func (m StreamMessage) ToMap() map[string]any {
	return map[string]any{
		"event":     m.Event,
		"user_id":   m.UserID,
		"chat_id":   m.ChatID,
		"content":   m.Content,
		"timestamp": m.Timestamp.UTC().Format(time.RFC3339Nano),
	}
}

func ParseStreamMessage(ackID string, data map[string]any) (StreamMessage, error) {
	var result StreamMessage

	result.AckID = ackID

	if chatID, ok := data[chatIDMapKey].(string); ok {
		result.ChatID = chatID
	} else {
		return StreamMessage{}, fmt.Errorf("invalid or missing chat_id")
	}

	if userID, ok := data[userIDMapKey].(string); ok {
		result.UserID = userID
	} else {
		return StreamMessage{}, fmt.Errorf("invalid or missing user_id")
	}

	if content, ok := data[contentMapKey].(string); ok {
		result.Content = content
	} else {
		return StreamMessage{}, fmt.Errorf("invalid or missing content")
	}

	if event, ok := data[eventMapKey].(string); ok {
		result.Event = event
	} else {
		return StreamMessage{}, fmt.Errorf("invalid or missing event")
	}

	if tsStr, ok := data[timestampKey].(string); ok {
		ts, err := time.Parse(time.RFC3339Nano, tsStr)
		if err != nil {
			return StreamMessage{}, fmt.Errorf("invalid timestamp format: %w", err)
		}
		result.Timestamp = ts
	} else {
		return StreamMessage{}, fmt.Errorf("invalid or missing timestamp")
	}

	return result, nil
}
