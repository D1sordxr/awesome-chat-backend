package chat

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/D1sordxr/awesome-chat-proto/gen/go/awesomechat/v1"

	"awesome-chat/internal/application/chat/query"
	"awesome-chat/internal/domain/core/chat/entity"
	"awesome-chat/internal/transport/grpc/handler"
)

func (h *HandlerImpl) ListChatPreviews(
	ctx context.Context,
	_ *v1.ListChatPreviewsRequest,
) (*v1.ListChatPreviewsResponse, error) {
	callerID, err := handler.CallerID(ctx)
	if err != nil {
		return nil, err
	}

	out, err := h.uc.Previews(ctx, query.PreviewsIn{UserID: callerID})
	if err != nil {
		return nil, err
	}

	previews := make([]*v1.ChatPreview, 0, len(out.Chats))
	for _, preview := range out.Chats {
		previews = append(previews, toProtoChatPreview(preview))
	}

	return &v1.ListChatPreviewsResponse{
		Chats: previews,
		Page:  &v1.PageResponse{Count: handler.Int32(len(previews))},
	}, nil
}

func toProtoChatPreview(preview entity.ChatPreview) *v1.ChatPreview {
	participants := make([]*v1.Participant, 0, len(preview.Participants))
	for _, participant := range preview.Participants {
		participants = append(participants, &v1.Participant{
			UserId:    participant.UserID.String(),
			Username:  participant.Username,
			AvatarUrl: participant.AvatarURL,
		})
	}

	chat := &v1.ChatPreview{
		ChatId:       preview.ChatID.String(),
		Name:         preview.Name,
		UnreadCount:  handler.Int32(preview.UnreadCount),
		AvatarUrl:    preview.AvatarURL,
		Participants: participants,
	}

	if preview.LastMessage.SenderID != uuid.Nil {
		chat.LastMessage = &v1.Message{
			MessageId: preview.LastMessage.ID,
			ChatId:    preview.ChatID.String(),
			SenderId:  preview.LastMessage.SenderID.String(),
			Type:      v1.MessageType_MESSAGE_TYPE_TEXT,
			Content:   preview.LastMessage.Text,
			CreatedAt: timestamppb.New(preview.LastMessage.Timestamp),
		}
	}

	return chat
}
