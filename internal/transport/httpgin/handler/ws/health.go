package ws

import "github.com/gin-gonic/gin"

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "ok"})
}

func (h *HealthHandler) RegisterRoutes(router gin.IRouter) {
	router.GET("api/ws/health", h.Check)
}
