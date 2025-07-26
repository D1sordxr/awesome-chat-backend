package health

import "github.com/gofiber/fiber/v2"

type Handler struct{}

func (h *Handler) Check(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/health", h.Check)
}
