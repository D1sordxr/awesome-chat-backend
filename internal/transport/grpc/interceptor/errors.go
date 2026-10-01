package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/transport/grpc/errors"
	"awesome-chat/internal/transport/grpc/logfields"
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

		code := status.Code(mapped)

		logByCode(e.log, code, "Request error", append(logfields.From(ctx),
			"method", info.FullMethod,
			"code", code.String(),
			"error", err.Error(),
		)...)

		return nil, mapped
	}
}
