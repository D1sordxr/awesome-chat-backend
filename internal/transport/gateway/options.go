package gateway

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/infrastructure/config/components/cookie"
	"awesome-chat/internal/transport/grpc/interceptor"
)

const (
	authorizationKey = "authorization"
	bearerPrefix     = "Bearer "
	expiredMaxAge    = -1
	requestIDHeader  = "X-Request-Id"
)

func newMux(cfg cookie.Config) *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			UseProtoNames:   true,
			EmitUnpopulated: true,
			DiscardUnknown:  true,
		}),
		runtime.WithMetadata(annotator(cfg.Name)),
		runtime.WithOutgoingHeaderMatcher(outgoingHeaderMatcher),
		runtime.WithForwardResponseOption(authCookieForwarder(cfg)),
	)
}

func annotator(cookieName string) func(context.Context, *http.Request) metadata.MD {
	return func(_ context.Context, r *http.Request) metadata.MD {
		md := metadata.MD{}

		if requestID := r.Header.Get(requestIDHeader); requestID != "" {
			md.Set(interceptor.RequestIDKey, requestID)
		}

		if r.Header.Get(authorizationKey) != "" {
			return md
		}

		authCookie, err := r.Cookie(cookieName)
		if err != nil || authCookie.Value == "" {
			return md
		}

		md.Set(authorizationKey, bearerPrefix+authCookie.Value)

		return md
	}
}

func outgoingHeaderMatcher(key string) (string, bool) {
	if strings.EqualFold(key, interceptor.RequestIDKey) {
		return requestIDHeader, true
	}

	return runtime.MetadataHeaderPrefix + key, true
}

func authCookieForwarder(cfg cookie.Config) func(context.Context, http.ResponseWriter, proto.Message) error {
	return func(_ context.Context, w http.ResponseWriter, message proto.Message) error {
		switch response := message.(type) {
		case *v1.LoginResponse:
			http.SetCookie(w, newAuthCookie(cfg, response.GetToken(), int(cfg.TTL/time.Second)))
		case *v1.LogoutResponse:
			http.SetCookie(w, newAuthCookie(cfg, "", expiredMaxAge))
		}

		return nil
	}
}

func newAuthCookie(cfg cookie.Config, value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     cfg.Name,
		Value:    value,
		Path:     cfg.Path,
		Domain:   cfg.Domain,
		MaxAge:   maxAge,
		Secure:   cfg.Secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}
