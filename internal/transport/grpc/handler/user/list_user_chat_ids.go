package user

import (
	"context"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/application/user/query"
	"awesome-chat/internal/pkg/validate"
	"awesome-chat/internal/transport/grpc/handler"
)

func (h *HandlerImpl) ListUserChatIds(
	ctx context.Context,
	req *v1.ListUserChatIdsRequest,
) (*v1.ListUserChatIdsResponse, error) {
	if err := validate.Proto(req); err != nil {
		return nil, err
	}

	userID, err := handler.ParseUUID(req.GetUserId(), "user_id")
	if err != nil {
		return nil, err
	}

	out, err := h.uc.ChatIDs(ctx, query.ChatIDsIn{UserID: userID})
	if err != nil {
		return nil, err
	}

	chatIDs := make([]string, 0, len(out.ChatIDs))
	for _, chatID := range out.ChatIDs {
		chatIDs = append(chatIDs, chatID.String())
	}

	return &v1.ListUserChatIdsResponse{ChatIds: chatIDs}, nil
}
