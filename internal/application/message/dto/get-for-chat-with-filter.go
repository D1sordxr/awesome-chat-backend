package dto

type (
	GetForChatWithFilterRequest struct {
		ChatID string `json:"chat_id"`
		Limit  int    `json:"limit,omitempty"`  // 100
		Offset int    `json:"offset,omitempty"` // pagination
		// or
		Cursor int `json:"cursor,omitempty"` // last message ID (cursor pagination)
	}
	GetForChatWithFilterResponse struct {
		AllMessages []FilteredMessage `json:"messages"`
		Count       int               `json:"count"`
		LastCursor  int               `json:"last_cursor"`
	}
	FilteredMessage struct {
		ID        int    `json:"id"`
		Text      string `json:"text"`
		SenderID  string `json:"sender_id"`
		Timestamp string `json:"timestamp"`
	}
)
