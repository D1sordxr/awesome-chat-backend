package entity

import (
	"awesome-chat/internal/domain/core/user/vo"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID uuid.UUID

	Email    vo.Email
	Password []byte
	Username string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(
	userID uuid.UUID,
	username string,
	email vo.Email,
	password []byte,
) User {
	return User{
		UserID:    userID,
		Email:     email,
		Password:  password,
		Username:  username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type UserRaw struct {
	UserID uuid.UUID

	Email    string
	Password []byte
	Username string

	CreatedAt time.Time
	UpdatedAt time.Time
}
