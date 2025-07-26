package middleware

import (
	"awesome-chat/internal/domain/core/user/ports"
	"awesome-chat/internal/domain/core/user/vo"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	authenticator ports.GetUserChatIDsStore
	timeout       time.Duration
}

func NewAuth(authenticator ports.GetUserChatIDsStore) *Auth {
	return &Auth{
		authenticator: authenticator,
		timeout:       time.Second * 2,
	}
}

func (a *Auth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if _, err := uuid.Parse(userID); err != nil {
			http.Error(w, "invalid user_id format", http.StatusBadRequest)
			return
		}

		authCtx, cancel := context.WithTimeout(r.Context(), a.timeout)
		defer cancel()

		chatIDs, err := a.authenticator.Execute(authCtx, userID)
		if err != nil {
			http.Error(w, "failed to get user chat", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), vo.UserIDKey, userID)
		ctx = context.WithValue(ctx, vo.ChatIDsKey, chatIDs)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
