package user

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/application/user/query"
	"awesome-chat/internal/domain/core/user/entity"
	"awesome-chat/internal/transport/grpc/handler"
)

func (h *HandlerImpl) ListUsers(
	ctx context.Context,
	_ *v1.ListUsersRequest,
) (*v1.ListUsersResponse, error) {
	out, err := h.uc.List(ctx, query.ListIn{})
	if err != nil {
		return nil, err
	}

	users := make([]*v1.User, 0, len(out.Users))
	for _, user := range out.Users {
		users = append(users, toProtoUser(user))
	}

	return &v1.ListUsersResponse{
		Users: users,
		Page:  &v1.PageResponse{Count: handler.Int32(len(users))},
	}, nil
}

func toProtoUser(user entity.User) *v1.User {
	return &v1.User{
		UserId:    user.UserID.String(),
		Username:  user.Username,
		Email:     user.Email.String(),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
