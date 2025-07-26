package ws

import (
	"awesome-chat/internal/domain/core/shared/ports/ws"
	"awesome-chat/internal/domain/core/user/ports"
	"awesome-chat/internal/domain/core/user/vo"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type UpgradeHandler struct {
	manager ws.ConnManager
	authMid ports.GinAuthMiddleware
}

func NewUpgrade(
	manager ws.ConnManager,
	authMid ports.GinAuthMiddleware,
) *UpgradeHandler {
	return &UpgradeHandler{
		manager: manager,
		authMid: authMid,
	}
}

func (h *UpgradeHandler) Do(ctx *gin.Context) {
	if !websocket.IsWebSocketUpgrade(ctx.Request) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "websocket upgrade required"})
		return
	}

	userID, exists := ctx.Get(string(vo.UserIDKey))
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user id not found"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	chatIDs, exists := ctx.Get(string(vo.ChatIDsKey))
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "chat ids not found"})
		return
	}

	chatIDsSlice, ok := chatIDs.([]string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid chat ids"})
		return
	}

	requestedID := ctx.Param("id")
	if requestedID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id parameter required"})
		return
	}

	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "UpgradeHandler, Connection, Sec-WebSocket-Version, Sec-WebSocket-Key")

	if err := h.manager.HandleWebSocket(
		ctx.Writer,
		ctx.Request,
		ctx.Request.Header,
		userIDStr,
		chatIDsSlice,
	); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to upgrade to websocket",
			"details": err.Error(),
		})
	}
}

func (h *UpgradeHandler) RegisterRoutes(router gin.IRouter) {
	router.GET("/api/ws/upgrade/:id", h.authMid.Auth(), h.Do)
}
