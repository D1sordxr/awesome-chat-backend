package interceptor

import (
	"awesome-chat/internal/transport/grpc/errors"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"awesome-chat/internal/domain/app/ports"
)

type Error struct {
	log ports.Logger
}

func NewError(log ports.Logger) *Error {
	return &Error{log: log}
}

func (e *Error) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		mapped := errors.MapError(err)

		logByCode(e.log, status.Code(mapped), "Request error",
			"method", info.FullMethod,
			"code", status.Code(mapped).String(),
			"error", err.Error(),
		)

		return nil, mapped
	}
}
