package dto

type (
	User struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	Users []User
)
type UserWS struct {
	UserID  string   `json:"user_id"`
	ChatIDs []string `json:"chat_ids"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
}

type UserResponse struct {
	ID        string `json:"id,omitempty"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type (
	UserID  string
	ChatIDs []string
)
