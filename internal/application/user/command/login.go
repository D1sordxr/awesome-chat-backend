package command

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"awesome-chat/internal/application/user/port"
	"awesome-chat/internal/domain/core/user/vo"
)

type Login struct {
	provider port.Provider
	tokens   port.TokenCreator
}

func NewLogin(provider port.Provider, tokens port.TokenCreator) *Login {
	return &Login{provider: provider, tokens: tokens}
}

type LoginIn struct {
	Email    string
	Password string
}

type LoginOut struct {
	UserID   string
	Username string
	Token    string
}

func (c *Login) Handle(ctx context.Context, in LoginIn) (LoginOut, error) {
	email, err := vo.NewEmail(in.Email)
	if err != nil {
		return LoginOut{}, err
	}

	user, err := c.provider.ByEmail(ctx, email.String())
	if err != nil {
		return LoginOut{}, err
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(in.Password)); err != nil {
		return LoginOut{}, err
	}

	token, err := c.tokens.Do(user)
	if err != nil {
		return LoginOut{}, err
	}

	return LoginOut{
		UserID:   user.UserID.String(),
		Username: user.Username,
		Token:    token,
	}, nil
}
