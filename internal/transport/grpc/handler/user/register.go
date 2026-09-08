package user

import (
	"context"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/application/user/command"
	"awesome-chat/internal/pkg/validate"
)

func (h *HandlerImpl) Register(
	ctx context.Context,
	req *v1.RegisterRequest,
) (*v1.RegisterResponse, error) {
	if err := validate.Proto(req); err != nil {
		return nil, err
	}

	out, err := h.uc.Register(ctx, command.RegisterIn{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.RegisterResponse{UserId: out.UserID.String()}, nil
}
