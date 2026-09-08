package gateway

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/infrastructure/config/components/cookie"
)

const (
	authorizationKey = "authorization"
	bearerPrefix     = "Bearer "
	expiredMaxAge    = -1
)

func newMux(cfg cookie.Config) *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
		}),
		runtime.WithMetadata(cookieAnnotator(cfg.Name)),
		runtime.WithForwardResponseOption(authCookieForwarder(cfg)),
	)
}

func cookieAnnotator(cookieName string) func(context.Context, *http.Request) metadata.MD {
	return func(_ context.Context, r *http.Request) metadata.MD {
		if r.Header.Get(authorizationKey) != "" {
			return nil
		}

		authCookie, err := r.Cookie(cookieName)
		if err != nil || authCookie.Value == "" {
			return nil
		}

		return metadata.Pairs(authorizationKey, bearerPrefix+authCookie.Value)
	}
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
