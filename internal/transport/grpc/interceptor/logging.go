package interceptor

import (
	"context"
	"log/slog"
	"time"

	logging "github.com/D1sordxr/packages/log"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const RequestIDKey = "x-request-id"

type Logger struct {
	log *slog.Logger
}

func NewLogger(log *slog.Logger) *Logger {
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

		requestID := requestID(ctx)
		ctx = logging.Inject(ctx, "request_id", requestID)

		_ = grpc.SetHeader(ctx, metadata.Pairs(RequestIDKey, requestID))

		resp, err := handler(ctx, req)

		code := status.Code(err)

		logByCode(ctx, l.log, code, "Request handled",
			"method", info.FullMethod,
			"code", code.String(),
			"duration", time.Since(started).String(),
		)

		return resp, err
	}
}

func requestID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if values := md.Get(RequestIDKey); len(values) > 0 && values[0] != "" {
			return values[0]
		}
	}

	return uuid.NewString()
}

func logByCode(ctx context.Context, log *slog.Logger, code codes.Code, msg string, fields ...any) {
	level := slog.LevelWarn

	switch code {
	case codes.OK:
		level = slog.LevelInfo
	case codes.Internal, codes.Unknown, codes.DataLoss, codes.Unavailable:
		level = slog.LevelError
	}

	log.Log(ctx, level, msg, fields...)
}
