package user

import (
	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/user/vo"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Auth struct {
	log ports.Logger
	uc  authJWTUseCase
}

func NewAuth(log ports.Logger, uc authJWTUseCase) *Auth {
	return &Auth{log: log, uc: uc}
}

func (h *Auth) authJWT(ctx *fiber.Ctx) error {
	start := time.Now()

	h.log.Info("AuthJWT handler started",
		"path", ctx.Path(),
		"method", ctx.Method(),
		"ip", ctx.IP(),
	)

	defer func() {
		h.log.Debug("AuthJWT handler completed",
			"duration", time.Since(start).String(),
		)
	}()

	tokenStr := ctx.Cookies("jwt")
	h.log.Debug("JWT token received",
		"token_length", len(tokenStr),
		"token_truncated", safeTruncate(tokenStr, 10),
	)

	if tokenStr == "" {
		h.log.Warn("JWT token missing in cookies")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token required",
		})
	}

	h.log.Debug("Calling auth use case")
	user, err := h.uc.Execute(ctx.Context(), vo.JWTToken(tokenStr))
	if err != nil {
		h.log.Error("Auth use case failed",
			"error", err,
			"stack", getCallerInfo(),
		)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	h.log.Info("User authenticated successfully",
		"user_id", user.UserID,
		"username", user.Username,
	)

	return ctx.JSON(fiber.Map{
		"username": user.Username,
		"id":       user.UserID,
	})
}

func (h *Auth) RegisterRoutes(router fiber.Router) {
	h.log.Info("Registering auth routes")
	router.Get("/user/jwt", h.authJWT)
	h.log.Info("Auth route registered",
		"path", "/user/jwt",
		"method", "GET",
	)
}

// Вспомогательные функции

func safeTruncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func getCallerInfo() string {
	// Реализация получения информации о вызове
	return "authJWT handler"
}
