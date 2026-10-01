package logfields

import (
	"context"
	"sync"
)

type contextKey struct{}

type holder struct {
	mu     sync.Mutex
	fields []any
}

func Inject(ctx context.Context, fields ...any) context.Context {
	return context.WithValue(ctx, contextKey{}, &holder{
		fields: append([]any(nil), fields...),
	})
}

func Set(ctx context.Context, key string, value any) {
	h, ok := ctx.Value(contextKey{}).(*holder)
	if !ok {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.fields = append(h.fields, key, value)
}

func From(ctx context.Context) []any {
	h, ok := ctx.Value(contextKey{}).(*holder)
	if !ok {
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	return append([]any(nil), h.fields...)
}
