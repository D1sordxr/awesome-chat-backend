package chat

import (
	"awesome-chat/internal/application/chat/dto"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type (
	createChatWithMembersUseCase interface {
		Execute(ctx context.Context, req dto.CreateChatRequest) (*dto.ChatResponse, error)
	}
	addUserUseCase interface {
		Execute(ctx context.Context, req dto.AddUserRequest) error
	}
	getUserChatPreviewUseCase interface {
		Execute(ctx context.Context, userID dto.UserID) (dto.GetUserChatPreviewResponse, error)
	}
	getChatAllMessagesUseCase interface {
		Execute(ctx context.Context, chatID dto.ChatID) (dto.AllMessages, error)
	}
)

type Handler struct {
	createUC             createChatWithMembersUseCase
	addUserUC            addUserUseCase
	getUserChatPreviewUC getUserChatPreviewUseCase
	getChatAllMessagesUC getChatAllMessagesUseCase
}

func NewChatHandler(
	createUC createChatWithMembersUseCase,
	addUserUC addUserUseCase,
	getUserChatPreviewUC getUserChatPreviewUseCase,
	getChatAllMessagesUC getChatAllMessagesUseCase,
) *Handler {
	return &Handler{
		createUC:             createUC,
		addUserUC:            addUserUC,
		getUserChatPreviewUC: getUserChatPreviewUC,
		getChatAllMessagesUC: getChatAllMessagesUC,
	}
}

func (h *Handler) CreateChatWithMembers(ctx *fiber.Ctx) error {
	var req dto.CreateChatRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	chat, err := h.createUC.Execute(ctx.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(chat)
}

func (h *Handler) AddUser(ctx *fiber.Ctx) error {
	var req dto.AddUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.addUserUC.Execute(ctx.Context(), req); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) getUserChatPreview(ctx *fiber.Ctx) error {
	reqCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user id required",
		})
	}

	resp, err := h.getUserChatPreviewUC.Execute(reqCtx, dto.UserID(id))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal server error",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"chat_previews": resp.ChatPreviews,
		"meta": fiber.Map{
			"count": len(resp.ChatPreviews),
		},
	})
}

func (h *Handler) getChatAllMessages(ctx *fiber.Ctx) error {
	reqCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	id := ctx.Params("chat_id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "chat id required",
		})
	}

	messages, err := h.getChatAllMessagesUC.Execute(reqCtx, dto.ChatID(id))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal server error",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(messages)
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/chat", h.CreateChatWithMembers)
	router.Post("/chat/add-user", h.AddUser)
	router.Get("/chat/:id", h.getUserChatPreview)
	router.Get("/chat/messages/:chat_id", h.getChatAllMessages)
}
