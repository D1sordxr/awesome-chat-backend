package ports

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type GetUserChatIDsStore interface {
	Execute(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type AuthMiddleware interface {
	Auth(http.Handler) http.Handler
}

type GinAuthMiddleware interface {
	Auth() gin.HandlerFunc
}
