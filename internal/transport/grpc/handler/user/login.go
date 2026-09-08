package user

import (
	"context"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/application/user/command"
	"awesome-chat/internal/pkg/validate"
)

func (h *HandlerImpl) Login(
	ctx context.Context,
	req *v1.LoginRequest,
) (*v1.LoginResponse, error) {
	if err := validate.Proto(req); err != nil {
		return nil, err
	}

	out, err := h.uc.Login(ctx, command.LoginIn{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.LoginResponse{
		UserId:   out.UserID,
		Username: out.Username,
		Token:    out.Token,
	}, nil
}
