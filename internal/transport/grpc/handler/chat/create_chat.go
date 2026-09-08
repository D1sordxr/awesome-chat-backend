package chat

import (
	"context"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/application/chat/command"
	"awesome-chat/internal/pkg/validate"
	"awesome-chat/internal/transport/grpc/handler"
)

func (h *HandlerImpl) CreateChat(
	ctx context.Context,
	req *v1.CreateChatRequest,
) (*v1.CreateChatResponse, error) {
	callerID, err := handler.CallerID(ctx)
	if err != nil {
		return nil, err
	}

	if err = validate.Proto(req); err != nil {
		return nil, err
	}

	out, err := h.uc.Create(ctx, command.CreateIn{
		OwnerID:   callerID,
		Name:      req.GetName(),
		MemberIDs: req.GetMemberIds(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateChatResponse{
		Chat: &v1.ChatPreview{
			ChatId: out.ChatID.String(),
			Name:   out.Name,
		},
	}, nil
}
