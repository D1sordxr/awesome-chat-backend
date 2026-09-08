package identity

import (
	"context"

	"github.com/google/uuid"
)

type contextKey struct{}

type Identity struct {
	UserID uuid.UUID
	Email  string
}

func With(ctx context.Context, id Identity) context.Context {
	return context.WithValue(ctx, contextKey{}, id)
}

func From(ctx context.Context) (Identity, bool) {
	id, ok := ctx.Value(contextKey{}).(Identity)

	return id, ok
}
