package interceptor

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"awesome-chat/internal/transport/grpc/errors"
)

type Error struct {
	log *slog.Logger
}

func NewError(log *slog.Logger) *Error {
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

		logByCode(ctx, e.log, code, "Request error",
			"method", info.FullMethod,
			"code", code.String(),
			"error", err.Error(),
		)

		return nil, mapped
	}
}
