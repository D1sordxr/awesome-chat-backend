package messages

import (
	"awesome-chat/internal/application/message/dto"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type (
	saveUseCase interface {
		Execute(ctx context.Context, msg dto.Message) error
	}
	getMessagesUseCase interface {
		Execute(ctx context.Context, req dto.GetRequest) (dto.Messages, error)
	}
	sendFastUseCase interface {
		Execute(ctx context.Context, payload []byte) error
	}
	sendUseCase interface {
		Execute(ctx context.Context, req dto.SendRequest) error
	}
	sendSyncUseCase interface {
		Execute(ctx context.Context, req dto.SendSyncRequest, raw []byte) error
	}
	getForChatWithFilterUseCase interface {
		Execute(ctx context.Context, req dto.GetForChatWithFilterRequest) (dto.GetForChatWithFilterResponse, error)
	}
)

type Handler struct {
	getMessagesUC          getMessagesUseCase
	saveUC                 saveUseCase
	sendUC                 sendUseCase
	sendFastUC             sendFastUseCase
	sendSyncUC             sendSyncUseCase
	getForChatWithFilterUC getForChatWithFilterUseCase
}

func NewMessageHandler(
	getMessagesUC getMessagesUseCase,
	saveUC saveUseCase,
	sendUC sendUseCase,
	sendFastUC sendFastUseCase,
	sendSyncUC sendSyncUseCase,
	getForChatWithFilterUC getForChatWithFilterUseCase,
) *Handler {
	return &Handler{
		sendSyncUC:             sendSyncUC,
		saveUC:                 saveUC,
		sendUC:                 sendUC,
		sendFastUC:             sendFastUC,
		getMessagesUC:          getMessagesUC,
		getForChatWithFilterUC: getForChatWithFilterUC,
	}
}

func (h *Handler) Save(ctx *fiber.Ctx) error {
	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 3*time.Second)
	defer cancel()

	var msg dto.Message
	if err := ctx.BodyParser(&msg); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.saveUC.Execute(timeoutCtx, msg); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) SendFast(ctx *fiber.Ctx) error {
	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cancel()

	if err := h.sendFastUC.Execute(timeoutCtx, ctx.Body()); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusAccepted)
}

func (h *Handler) Send(ctx *fiber.Ctx) error {
	var data dto.SendRequest
	err := ctx.BodyParser(&data)
	if err != nil {
		return err
	}

	if err = h.sendUC.Execute(ctx.Context(), data); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusAccepted)
}

func (h *Handler) SendSync(ctx *fiber.Ctx) error {
	var data dto.SendSyncRequest
	err := ctx.BodyParser(&data)
	if err != nil {
		return err
	}

	if err = h.sendSyncUC.Execute(ctx.Context(), data, ctx.Body()); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusAccepted)
}

func (h *Handler) GetMessages(ctx *fiber.Ctx) error {
	chatID := ctx.Query("chat_id")
	if chatID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "chat_id is required")
	}

	limit := ctx.QueryInt("limit", 100)

	offset := ctx.QueryInt("offset", 0)

	cursor := ctx.Query("cursor", "")

	messages, err := h.getMessagesUC.Execute(ctx.Context(), dto.GetRequest{
		ChatID: chatID,
		Limit:  limit,
		Offset: offset,
		Cursor: cursor,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(messages)
}

func (h *Handler) getForChatWithFilter(ctx *fiber.Ctx) error {
	chatID := ctx.Query("chat_id")
	if chatID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "chat_id is required")
	}

	limit := ctx.QueryInt("limit", 100)

	offset := ctx.QueryInt("offset", 0)

	cursor := ctx.QueryInt("cursor", 0)

	resp, err := h.getForChatWithFilterUC.Execute(ctx.Context(), dto.GetForChatWithFilterRequest{
		ChatID: chatID,
		Limit:  limit,
		Offset: offset,
		Cursor: cursor,
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"messages": resp.AllMessages,
	})

}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/message/save", h.Save)
	router.Post("/message/send", h.Send)
	router.Post("/message/send-fast", h.SendFast)
	router.Post("/message/send-sync", h.SendSync)
	router.Get("/message", h.GetMessages)
	router.Get("/message/get-for-chat-with-filter", h.SendFast)
}
