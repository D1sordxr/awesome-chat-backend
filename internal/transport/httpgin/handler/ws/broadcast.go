package ws

import (
	"awesome-chat/internal/domain/core/shared/ports/ws"
	"awesome-chat/internal/infrastructure/ws/chathub"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type BroadcastHandler struct {
	uc ws.MessageBroadcaster
}

func NewBroadcastHandler(uc ws.MessageBroadcaster) *BroadcastHandler {
	return &BroadcastHandler{
		uc: uc,
	}
}

func (h *BroadcastHandler) BroadcastMessage(ctx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var message chathub.Message
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.Broadcast(reqCtx, message); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *BroadcastHandler) RegisterRoutes(router gin.IRouter) {
	router.POST("/api/ws/broadcast", h.BroadcastMessage)
}
