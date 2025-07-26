package dto

type CreateChatRequest struct {
	Name      string   `json:"name"`
	MemberIDs []string `json:"member_ids"`
}

type ChatResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type AddUserRequest struct {
	ChatID string `json:"chat_id" validate:"required"`
	UserID string `json:"user_id" validate:"required"`
}
