package chat

import (
	"context"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/application/chat/command"
	"awesome-chat/internal/pkg/validate"
	"awesome-chat/internal/transport/grpc/handler"
)

func (h *HandlerImpl) AddMember(
	ctx context.Context,
	req *v1.AddMemberRequest,
) (*v1.AddMemberResponse, error) {
	callerID, err := handler.CallerID(ctx)
	if err != nil {
		return nil, err
	}

	if err = validate.Proto(req); err != nil {
		return nil, err
	}

	chatID, err := handler.ParseUUID(req.GetChatId(), "chat_id")
	if err != nil {
		return nil, err
	}

	userID, err := handler.ParseUUID(req.GetUserId(), "user_id")
	if err != nil {
		return nil, err
	}

	if _, err = h.uc.AddMember(ctx, command.AddMemberIn{
		CallerID: callerID,
		ChatID:   chatID,
		UserID:   userID,
	}); err != nil {
		return nil, err
	}

	return &v1.AddMemberResponse{}, nil
}
