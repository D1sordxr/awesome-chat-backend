package sendMessage

import (
	"awesome-chat/internal/application/message/dto"
	"awesome-chat/internal/domain/core/message/ports/usecases"
	"awesome-chat/internal/infrastructure/ws/chathub"
	"awesome-chat/internal/infrastructure/ws/chathub/consts"
	"awesome-chat/internal/infrastructure/ws/chathub/transport"
	"context"
	"encoding/json"
	"fmt"
)

type Handler struct {
	opType consts.OperationType
	uc     usecases.MessageBroadcastWithPub
}

func New(uc usecases.MessageBroadcastWithPub) *Handler {
	return &Handler{
		opType: consts.SendMessage,
		uc:     uc,
	}
}

func (h *Handler) Handle(ctx context.Context, body json.RawMessage) chathub.OperationResponse {
	var req dto.BroadcastWithPubRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return chathub.ErrorResponse(h.opType.String(), fmt.Errorf("invalid body format: %w", err))
	}

	if err := h.uc.Execute(ctx, req); err != nil {
		return chathub.ErrorResponse(h.opType.String(), fmt.Errorf("send message error: %w", err))
	}

	return chathub.SuccessResponse(h.opType.String(), map[string]string{
		"message": "success",
	})
}

func (h *Handler) Register(handlerStore transport.HandlerStore) {
	handlerStore[h.opType] = h
}
