package user

import (
	"context"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/application/user/query"
	"awesome-chat/internal/transport/grpc/handler"
)

func (h *HandlerImpl) GetCurrentUser(
	ctx context.Context,
	_ *v1.GetCurrentUserRequest,
) (*v1.GetCurrentUserResponse, error) {
	callerID, err := handler.CallerID(ctx)
	if err != nil {
		return nil, err
	}

	out, err := h.uc.Read(ctx, query.ReadIn{UserID: callerID})
	if err != nil {
		return nil, err
	}

	return &v1.GetCurrentUserResponse{User: toProtoUser(out.User)}, nil
}
