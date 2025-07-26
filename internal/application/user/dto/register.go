package dto

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type RegisterUserResponse struct {
	UserID string `json:"user_id"`
}
