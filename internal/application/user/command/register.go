package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"awesome-chat/internal/application/user/port"
	"awesome-chat/internal/domain/core/user/entity"
	"awesome-chat/internal/domain/core/user/vo"
)

type Register struct {
	repo port.Repository
}

func NewRegister(repo port.Repository) *Register {
	return &Register{repo: repo}
}

type RegisterIn struct {
	Username string
	Email    string
	Password string
}

type RegisterOut struct {
	UserID uuid.UUID
}

func (c *Register) Handle(ctx context.Context, in RegisterIn) (RegisterOut, error) {
	email, err := vo.NewEmail(in.Email)
	if err != nil {
		return RegisterOut{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return RegisterOut{}, fmt.Errorf("hash password: %w", err)
	}

	userID := uuid.New()

	if err = c.repo.Save(ctx, entity.NewUser(userID, in.Username, email, hash)); err != nil {
		return RegisterOut{}, err
	}

	return RegisterOut{UserID: userID}, nil
}
