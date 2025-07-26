package middleware

import (
	"awesome-chat/internal/domain/core/user/ports"
	"awesome-chat/internal/domain/core/user/vo"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"time"
)

const (
	host = "api:8080"
	path = "/user/chat-ids/"
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

func (a *Auth) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := ctx.Cookie("user_id")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "no user id cookie"})
			ctx.Abort()
			return
		}
		if err = uuid.Validate(userID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id cookie"})
			ctx.Abort()
			return
		}

		authCtx, cancel := context.WithTimeout(ctx.Request.Context(), a.timeout)
		defer cancel()

		chatIDs, err := a.authenticator.Execute(authCtx, uuid.New())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		nextCtx := context.WithValue(ctx.Request.Context(), vo.UserIDKey, userID)
		nextCtx = context.WithValue(nextCtx, vo.ChatIDsKey, chatIDs)
		ctx.Request = ctx.Request.WithContext(nextCtx)

		ctx.Next()
	}
}

type AuthAsClient struct {
	httpClient *http.Client
	timeout    time.Duration
}

func NewAuthAsClient() *AuthAsClient {
	return &AuthAsClient{
		httpClient: http.DefaultClient,
		timeout:    time.Second * 2,
	}
}

func (a *AuthAsClient) Auth() gin.HandlerFunc { // TODO: add and change localhost -> cfg.authHost
	return func(ctx *gin.Context) {
		userID := ctx.Param("id")

		if err := uuid.Validate(userID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id cookie"})
			ctx.Abort()
			return
		}

		authURL := url.URL{
			Scheme: "http",
			Host:   host,
			Path:   path + userID,
		}

		req, err := http.NewRequest(
			http.MethodGet,
			authURL.String(),
			nil,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			ctx.Abort()
			return
		}

		authCtx, cancel := context.WithTimeout(ctx.Request.Context(), a.timeout)
		defer cancel()
		req = req.WithContext(authCtx)

		resp, err := a.httpClient.Do(req)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "auth service unavailable: " + err.Error()})
			ctx.Abort()
			return
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			ctx.JSON(resp.StatusCode, gin.H{"error": "auth failed"})
			ctx.Abort()
			return
		}

		var authResp struct {
			ChatIDs []string `json:"chat_ids"`
		}
		if err = json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode response"})
			ctx.Abort()
			return
		}

		ctx.Set(string(vo.UserIDKey), userID)
		ctx.Set(string(vo.ChatIDsKey), authResp.ChatIDs)
		ctx.Next()
	}
}
