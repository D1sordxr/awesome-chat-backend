package ports

import (
	"awesome-chat/internal/domain/core/user/entity"
	"awesome-chat/internal/domain/core/user/vo"
)

type TokenCreator interface {
	Do(user entity.User) (string, error)
}

type TokenParser interface {
	Do(tokenStr vo.JWTToken) (vo.IDClaims, vo.EmailClaims, error)
}
