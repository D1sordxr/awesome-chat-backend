package transport

import (
	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/infrastructure/ws/chathub"
	"awesome-chat/internal/infrastructure/ws/chathub/consts"
	"awesome-chat/internal/infrastructure/ws/chathub/errors"
	"context"
	"encoding/json"
	"fmt"
)

type (
	Handler interface {
		Handle(ctx context.Context, body json.RawMessage) chathub.OperationResponse
		Register(handlerStore HandlerStore)
	}
	HandlerStore     map[consts.OperationType]Handler
	OperationHandler struct {
		log          ports.Logger
		handlerStore HandlerStore
	}
)

func NewOperationHandler(
	log ports.Logger,
	handlers ...Handler,
) *OperationHandler {
	validHandlers := make(map[consts.OperationType]Handler, len(handlers))

	for _, h := range handlers {
		h.Register(validHandlers)
	}

	return &OperationHandler{
		log:          log,
		handlerStore: validHandlers,
	}
}

func (o *OperationHandler) Handle(
	ctx context.Context,
	data []byte,
) chathub.OperationResponse {
	const op = "ws.chathub.OperationHandler"

	var opDTO OperationHeader
	if err := json.Unmarshal(data, &opDTO); err != nil {
		return chathub.ErrorResponse(func() string {
			if opDTO.Operation != "" {
				return opDTO.Operation
			} else {
				return "unknown"
			}
		}(), fmt.Errorf("%s: %w", op, errors.ErrInvalidOpFormat))
	}

	handler, exists := o.handlerStore[consts.OperationType(opDTO.Operation)]
	if !exists {
		return chathub.ErrorResponse(opDTO.Operation,
			fmt.Errorf("%s: %w", op, errors.ErrInvalidOpFormat),
		)
	}

	return handler.Handle(ctx, opDTO.Body)
}
