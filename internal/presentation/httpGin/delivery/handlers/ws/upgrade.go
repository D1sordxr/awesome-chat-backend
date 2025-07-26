package ws

import (
	"awesome-chat/internal/domain/core/shared/ports/ws"
	"awesome-chat/internal/domain/core/user/ports"
	"awesome-chat/internal/domain/core/user/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpgradeHandler struct {
	manager ws.Upgrader
	authMid ports.GinAuthMiddleware
}

func NewUpgradeHandler(
	manager ws.Upgrader,
	authMid ports.GinAuthMiddleware,
) *UpgradeHandler {
	return &UpgradeHandler{
		manager: manager,
		authMid: authMid,
	}
}

func (h *UpgradeHandler) HandleWebSocket(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user id required"})
		return
	}

	var chatIDs []string
	if err := ctx.ShouldBind(&chatIDs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "chat ids required"})
		return
	}

	if err := h.manager.HandleWebSocket(
		ctx.Request.Context(),
		ctx.Writer,
		ctx.Request,
		nil,
		userID,
		chatIDs...,
	); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to upgrade to websocket",
			"details": err.Error(),
		})
		return
	}
}

func (h *UpgradeHandler) Do(ctx *gin.Context) {
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

	if err := h.manager.HandleWebSocket(
		ctx.Request.Context(),
		ctx.Writer,
		ctx.Request,
		nil,
		userIDStr,
		chatIDsSlice...,
	); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to upgrade to websocket",
			"details": err.Error(),
		})
		return
	}
}

func (h *UpgradeHandler) Test(ctx *gin.Context) {
	var (
		userID  = "test_user_id"
		chatIDs = []string{"test_chat_id", "yo"}
	)

	if err := h.manager.HandleWebSocket(
		ctx.Request.Context(),
		ctx.Writer,
		ctx.Request,
		nil,
		userID,
		chatIDs...,
	); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to upgrade to websocket",
		})
		return
	}
}

func (h *UpgradeHandler) TestPostman(ctx *gin.Context) {
	var (
		userID  = "test_postman_user_id"
		chatIDs = []string{"test_chat_id", "yo"}
	)

	if err := h.manager.HandleWebSocket(
		ctx.Request.Context(),
		ctx.Writer,
		ctx.Request,
		nil,
		userID,
		chatIDs...,
	); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to upgrade to websocket",
		})
		return
	}
}

func (h *UpgradeHandler) RegisterRoutes(router gin.IRouter) {
	router.GET("/ws/omitted/:id", h.HandleWebSocket)

	router.GET("/ws/:id", h.authMid.Auth(), h.Do)
	router.GET("/ws/test/", h.Test)
	router.GET("/ws/test/postman", h.TestPostman)
}
