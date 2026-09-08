package user

import (
	"context"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"
)

func (h *HandlerImpl) Logout(
	_ context.Context,
	_ *v1.LogoutRequest,
) (*v1.LogoutResponse, error) {
	return &v1.LogoutResponse{}, nil
}
