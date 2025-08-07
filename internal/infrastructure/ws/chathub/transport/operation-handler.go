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

	o.log.Debug("handling operation", "op", op)

	var opDTO OperationDTO
	if err := json.Unmarshal(data, &opDTO); err != nil {
		return chathub.ErrorResponse(func() string {
			if opDTO.Operation != "" {
				return opDTO.Operation
			} else {
				return "unknown"
			}
		}(), fmt.Errorf("%s: %w", op, errors.ErrInvalidOpFormat))
	}

	o.log.Debug("operation parsed",
		"op", op,
		"struct", opDTO,
		"opDTO.ID", opDTO.ID,
		"opDTO.Operation", opDTO.Operation,
	)

	handler, exists := o.handlerStore[consts.OperationType(opDTO.Operation)]
	if !exists {
		return chathub.ErrorResponse(opDTO.Operation,
			fmt.Errorf("%s: %w", op, errors.ErrInvalidOpFormat),
		)
	}

	resp := handler.Handle(ctx, opDTO.Body)
	return resp

}
