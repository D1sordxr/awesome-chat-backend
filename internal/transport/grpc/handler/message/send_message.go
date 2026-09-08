package message

import (
	"context"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/application/message/command"
	"awesome-chat/internal/pkg/validate"
	"awesome-chat/internal/transport/grpc/handler"
)

func (h *HandlerImpl) SendMessage(
	ctx context.Context,
	req *v1.SendMessageRequest,
) (*v1.SendMessageResponse, error) {
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

	out, err := h.uc.Send(ctx, command.SendIn{
		SenderID: callerID,
		ChatID:   chatID,
		Content:  req.GetContent(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.SendMessageResponse{Message: toProtoMessage(out.Message)}, nil
}
