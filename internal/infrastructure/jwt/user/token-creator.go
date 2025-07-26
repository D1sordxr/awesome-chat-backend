package user

import (
	"awesome-chat/internal/domain/core/user/entity"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenCreatorImpl struct {
	secretKey string
}

func NewTokenCreator(secretKey string) *TokenCreatorImpl {
	return &TokenCreatorImpl{secretKey: secretKey}
}

func (t *TokenCreatorImpl) Do(user entity.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims[idClaims] = user.UserID
	claims[emailClaims] = user.Email.String()
	claims[expClaims] = time.Now().Add(expDuration).Unix()

	tokenString, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
