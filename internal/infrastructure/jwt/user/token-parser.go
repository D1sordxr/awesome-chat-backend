package user

import (
	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/user/vo"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type TokenParserImpl struct {
	log       ports.Logger
	secretKey string
}

func NewTokenParser(
	log ports.Logger,
	secretKey string,
) *TokenParserImpl {
	return &TokenParserImpl{
		log:       log,
		secretKey: secretKey,
	}
}

func (t *TokenParserImpl) Do(tokenStr vo.JWTToken) (vo.IDClaims, vo.EmailClaims, error) {
	token, err := jwt.Parse(string(tokenStr), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid token claims format")
	}

	id, ok := claims[idClaims].(string)
	if !ok {
		return "", "", fmt.Errorf("id claim missing or invalid")
	}
	t.log.Debug("ID claims received", "id", id)

	email, ok := claims[emailClaims].(string)
	if !ok {
		return "", "", fmt.Errorf("email claim missing or invalid")
	}
	t.log.Debug("ID claims received", "id", id)
	
	return vo.IDClaims(id), vo.EmailClaims(email), nil
}
