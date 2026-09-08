package interceptor

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/domain/core/user/vo"
	"awesome-chat/internal/transport/grpc/identity"
)

//go:generate mockgen -source=auth.go -destination=mocks/interceptors_mocks.go -package=mocks

const (
	authorizationKey = "authorization"
	bearerPrefix     = "bearer "
)

var PublicMethods = []string{
	v1.UserService_Register_FullMethodName,
	v1.UserService_Login_FullMethodName,
	v1.UserService_Logout_FullMethodName,
}

type tokenParser interface {
	Do(token vo.JWTToken) (vo.IDClaims, vo.EmailClaims, error)
}

type Auth struct {
	parser tokenParser
	public map[string]struct{}
}

func NewAuth(parser tokenParser, publicMethods ...string) *Auth {
	public := make(map[string]struct{}, len(publicMethods))
	for _, method := range publicMethods {
		public[method] = struct{}{}
	}

	return &Auth{parser: parser, public: public}
}

func (a *Auth) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if _, ok := a.public[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		caller, err := a.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(identity.With(ctx, caller), req)
	}
}

func (a *Auth) authenticate(ctx context.Context) (identity.Identity, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return identity.Identity{}, status.Error(codes.Unauthenticated, "authentication required")
	}

	values := md.Get(authorizationKey)
	if len(values) == 0 {
		return identity.Identity{}, status.Error(codes.Unauthenticated, "authentication required")
	}

	token := values[0]
	if len(token) >= len(bearerPrefix) && strings.EqualFold(token[:len(bearerPrefix)], bearerPrefix) {
		token = token[len(bearerPrefix):]
	}

	idClaim, emailClaim, err := a.parser.Do(vo.JWTToken(token))
	if err != nil {
		return identity.Identity{}, status.Error(codes.Unauthenticated, "invalid or expired token")
	}

	userID, err := uuid.Parse(string(idClaim))
	if err != nil {
		return identity.Identity{}, status.Error(codes.Unauthenticated, "invalid token subject")
	}

	return identity.Identity{UserID: userID, Email: string(emailClaim)}, nil
}
