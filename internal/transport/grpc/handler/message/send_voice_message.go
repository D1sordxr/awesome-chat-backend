package message

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/application/message/command"
	"awesome-chat/internal/domain/core/message/entity"
	"awesome-chat/internal/pkg/validate"
	"awesome-chat/internal/transport/grpc/handler"
)

func (h *HandlerImpl) SendVoiceMessage(
	ctx context.Context,
	req *v1.SendVoiceMessageRequest,
) (*v1.SendVoiceMessageResponse, error) {
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

	out, err := h.uc.SendVoice(ctx, command.SendVoiceIn{
		SenderID:        callerID,
		ChatID:          chatID,
		ObjectKey:       req.GetObjectKey(),
		DurationSeconds: int(req.GetDurationSeconds()),
	})
	if err != nil {
		return nil, err
	}

	message := toProtoMessage(out.Message)
	message.Type = v1.MessageType_MESSAGE_TYPE_VOICE
	message.Voice = &v1.Voice{
		ObjectKey:       out.Voice.ObjectKey,
		DurationSeconds: handler.Int32(out.Voice.DurationSeconds),
		Waveform:        out.Voice.Waveform,
	}

	return &v1.SendVoiceMessageResponse{Message: message}, nil
}

func toProtoMessage(message entity.Message) *v1.Message {
	return &v1.Message{
		MessageId: message.ID,
		ChatId:    message.ChatID.String(),
		SenderId:  message.UserID.String(),
		Type:      v1.MessageType_MESSAGE_TYPE_TEXT,
		Content:   message.Content,
		CreatedAt: timestamppb.New(message.Timestamp),
	}
}
