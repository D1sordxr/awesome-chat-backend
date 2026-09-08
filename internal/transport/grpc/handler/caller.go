package handler

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"awesome-chat/internal/transport/grpc/identity"
)

var ErrUnauthenticated = status.Error(codes.Unauthenticated, "authentication required")

func CallerID(ctx context.Context) (uuid.UUID, error) {
	caller, ok := identity.From(ctx)
	if !ok {
		return uuid.Nil, ErrUnauthenticated
	}

	return caller.UserID, nil
}
