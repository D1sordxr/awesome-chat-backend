package gateway

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/application/user/command"
	"awesome-chat/internal/application/user/query"
	"awesome-chat/internal/domain/core/user/entity"
	userVO "awesome-chat/internal/domain/core/user/vo"
	"awesome-chat/internal/infrastructure/config/components/cookie"
	httpCfg "awesome-chat/internal/infrastructure/config/components/http"
	"awesome-chat/internal/infrastructure/logger"
	userHandler "awesome-chat/internal/transport/grpc/handler/user"
	userMocks "awesome-chat/internal/transport/grpc/handler/user/mocks"
	"awesome-chat/internal/transport/grpc/interceptor"
	interceptorMocks "awesome-chat/internal/transport/grpc/interceptor/mocks"
)

const (
	testToken = "test-token"
	testEmail = "user@example.com"
)

type stubs struct {
	parser  *interceptorMocks.MocktokenParser
	useCase *userMocks.MockuseCase
}

func newGateway(t *testing.T) (*httptest.Server, stubs) {
	t.Helper()

	ctrl := gomock.NewController(t)
	s := stubs{
		parser:  interceptorMocks.NewMocktokenParser(ctrl),
		useCase: userMocks.NewMockuseCase(ctrl),
	}

	handler := userHandler.NewHandlerWithTracing(userHandler.NewHandler(s.useCase))

	auth := interceptor.NewAuth(s.parser, interceptor.PublicMethods...)
	grpcLogger := interceptor.NewLogger(logger.NewLogger())
	grpcErrors := interceptor.NewError(logger.NewLogger())
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcLogger.Unary(),
		grpcErrors.Unary(),
		auth.Unary(),
	))
	v1.RegisterUserServiceServer(grpcServer, handler)

	listener := bufconn.Listen(1024 * 1024)
	go func() { _ = grpcServer.Serve(listener) }()

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return listener.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial gRPC: %v", err)
	}

	gateway, err := NewServer(
		context.Background(),
		logger.NewLogger(),
		&httpCfg.Config{Port: "0", Timeout: time.Second, IdleTimeout: time.Second},
		Options{
			AllowedOrigins: []string{"http://localhost:3000"},
			Cookie:         cookie.Config{Name: "jwt", Path: "/", TTL: time.Hour},
		},
		conn,
		v1.RegisterUserServiceHandler,
	)
	if err != nil {
		t.Fatalf("build gateway: %v", err)
	}

	server := httptest.NewServer(gateway.Handler())
	t.Cleanup(func() {
		server.Close()
		_ = conn.Close()
		grpcServer.Stop()
	})

	return server, s
}

func TestLoginSetsHTTPOnlyCookie(t *testing.T) {
	t.Parallel()

	server, s := newGateway(t)
	userID := uuid.New()

	s.useCase.EXPECT().
		Login(gomock.Any(), command.LoginIn{Email: testEmail, Password: "secret"}).
		Return(command.LoginOut{
			UserID:   userID.String(),
			Username: "oleg",
			Token:    testToken,
		}, nil)

	resp, err := server.Client().Post(
		server.URL+"/v1/auth/login",
		"application/json",
		strings.NewReader(`{"email":"`+testEmail+`","password":"secret"}`),
	)
	if err != nil {
		t.Fatalf("login request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var authCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "jwt" {
			authCookie = c
		}
	}

	if authCookie == nil {
		t.Fatal("jwt cookie was not set")
	}
	if authCookie.Value != testToken {
		t.Errorf("cookie value = %q, want %q", authCookie.Value, testToken)
	}
	if !authCookie.HttpOnly {
		t.Error("cookie is not HttpOnly")
	}
}

func TestGetCurrentUserWithoutCookieIsUnauthenticated(t *testing.T) {
	t.Parallel()

	server, _ := newGateway(t)

	resp, err := server.Client().Get(server.URL + "/v1/auth/me")
	if err != nil {
		t.Fatalf("me request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestCookieIsForwardedAsAuthorizationMetadata(t *testing.T) {
	t.Parallel()

	server, s := newGateway(t)
	userID := uuid.New()

	s.parser.EXPECT().
		Do(userVO.JWTToken(testToken)).
		Return(userVO.IDClaims(userID.String()), userVO.EmailClaims(testEmail), nil)

	s.useCase.EXPECT().
		Read(gomock.Any(), query.ReadIn{UserID: userID}).
		Return(query.ReadOut{User: entity.User{
			UserID:   userID,
			Username: "oleg",
			Email:    userVO.Email(testEmail),
		}}, nil)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		server.URL+"/v1/auth/me",
		nil,
	)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.AddCookie(&http.Cookie{Name: "jwt", Value: testToken})

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("me request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var body struct {
		User struct {
			UserID string `json:"user_id"`
			Email  string `json:"email"`
		} `json:"user"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.User.UserID != userID.String() {
		t.Errorf("user_id = %q, want %q", body.User.UserID, userID.String())
	}
	if body.User.Email != testEmail {
		t.Errorf("email = %q, want %q", body.User.Email, testEmail)
	}
}

func TestLogoutExpiresCookie(t *testing.T) {
	t.Parallel()

	server, _ := newGateway(t)

	resp, err := server.Client().Post(server.URL+"/v1/auth/logout", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatalf("logout request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	for _, c := range resp.Cookies() {
		if c.Name == "jwt" {
			if c.MaxAge >= 0 {
				t.Errorf("MaxAge = %d, want negative", c.MaxAge)
			}

			return
		}
	}

	t.Fatal("jwt cookie was not cleared")
}

func TestRegisterRejectsInvalidEmail(t *testing.T) {
	t.Parallel()

	server, _ := newGateway(t)

	resp, err := server.Client().Post(
		server.URL+"/v1/auth/register",
		"application/json",
		strings.NewReader(`{"username":"oleg","email":"not-an-email","password":"secret-password"}`),
	)
	if err != nil {
		t.Fatalf("register request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}

	var body struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Details []struct {
			Type            string `json:"@type"`
			FieldViolations []struct {
				Field       string `json:"field"`
				Reason      string `json:"reason"`
				Description string `json:"description"`
			} `json:"field_violations"`
		} `json:"details"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(body.Details) == 0 || len(body.Details[0].FieldViolations) == 0 {
		t.Fatalf("no field violations in response: %+v", body)
	}

	violation := body.Details[0].FieldViolations[0]
	if violation.Field != "email" {
		t.Errorf("field = %q, want %q", violation.Field, "email")
	}
	if violation.Reason == "" {
		t.Error("violation has no rule id")
	}
}
