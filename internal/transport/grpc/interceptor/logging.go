package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"awesome-chat/internal/domain/app/ports"
)

type Logger struct {
	log ports.Logger
}

func NewLogger(log ports.Logger) *Logger {
	return &Logger{log: log}
}

func (l *Logger) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		started := time.Now()

		resp, err := handler(ctx, req)

		logByCode(l.log, status.Code(err), "Request handled",
			"method", info.FullMethod,
			"code", status.Code(err).String(),
			"duration", time.Since(started).String(),
		)

		return resp, err
	}
}

func logByCode(log ports.Logger, code codes.Code, msg string, fields ...any) {
	switch code {
	case codes.OK:
		log.Info(msg, fields...)
	case codes.Internal, codes.Unknown, codes.DataLoss, codes.Unavailable:
		log.Error(msg, fields...)
	default:
		log.Warn(msg, fields...)
	}
}
