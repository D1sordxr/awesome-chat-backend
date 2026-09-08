package message

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/application/message/command"
	"awesome-chat/internal/pkg/validate"
	"awesome-chat/internal/transport/grpc/handler"
)

func (h *HandlerImpl) CreateVoiceUpload(
	ctx context.Context,
	req *v1.CreateVoiceUploadRequest,
) (*v1.CreateVoiceUploadResponse, error) {
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

	out, err := h.uc.CreateVoiceUpload(ctx, command.CreateVoiceUploadIn{
		CallerID: callerID,
		ChatID:   chatID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateVoiceUploadResponse{
		ObjectKey: out.ObjectKey,
		UploadUrl: out.UploadURL,
		ExpiresAt: timestamppb.New(out.ExpiresAt),
	}, nil
}
