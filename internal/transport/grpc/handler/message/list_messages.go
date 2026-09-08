package message

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/application/message/query"
	"awesome-chat/internal/pkg/validate"
	"awesome-chat/internal/transport/grpc/handler"
)

func (h *HandlerImpl) ListMessages(
	ctx context.Context,
	req *v1.ListMessagesRequest,
) (*v1.ListMessagesResponse, error) {
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

	out, err := h.uc.List(ctx, query.ListIn{
		CallerID: callerID,
		ChatID:   chatID,
		Limit:    int(req.GetPage().GetLimit()),
		Cursor:   req.GetPage().GetCursor(),
	})
	if err != nil {
		return nil, err
	}

	messages := make([]*v1.Message, 0, len(out.Messages))
	for _, message := range out.Messages {
		messages = append(messages, &v1.Message{
			MessageId: message.ID,
			ChatId:    req.GetChatId(),
			SenderId:  message.SenderID.String(),
			Type:      v1.MessageType_MESSAGE_TYPE_TEXT,
			Content:   message.Text,
			CreatedAt: timestamppb.New(message.Timestamp),
		})
	}

	return &v1.ListMessagesResponse{
		Messages: messages,
		Page: &v1.PageResponse{
			Count:      handler.Int32(len(messages)),
			NextCursor: out.NextCursor,
			HasMore:    out.HasMore,
		},
	}, nil
}
