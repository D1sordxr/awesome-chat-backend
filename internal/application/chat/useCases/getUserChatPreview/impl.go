package getUserChatPreview

import (
	"awesome-chat/internal/application/chat/dto"
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/chat/ports"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type ChatGetUserChatPreviewUseCase struct {
	log   appPorts.Logger
	store ports.GetUserChatPreviewStore
}

func NewChatGetUserChatPreviewUseCase(
	log appPorts.Logger,
	store ports.GetUserChatPreviewStore,
) *ChatGetUserChatPreviewUseCase {
	return &ChatGetUserChatPreviewUseCase{
		log:   log,
		store: store,
	}
}

func (uc *ChatGetUserChatPreviewUseCase) Execute(
	ctx context.Context,
	userID dto.UserID,
) (
	dto.GetUserChatPreviewResponse,
	error,
) {
	const op = "ChatGetUserChatPreviewUseCase.SetupChatPreviews"
	withFields := func(args ...any) []any {
		return append([]any{"op", op, "user_id", string(userID)}, args...)
	}

	uc.log.Info("Attempting to get user chats", withFields()...)

	id, err := uuid.Parse(string(userID))
	if err != nil {
		uc.log.Error("Failed to parse user id", withFields("error", err.Error())...)
		return emptyWithErr(fmt.Errorf("%s: %w", op, err))
	}

	previews, err := uc.store.SetupChatPreviews(ctx, id)
	if err != nil {
		uc.log.Error("Failed to get user chats", withFields("error", err.Error())...)
		return emptyWithErr(fmt.Errorf("%s: %w", op, err))
	}

	chatIDs := make([]uuid.UUID, 0, len(previews))
	for _, p := range previews {
		chatIDs = append(chatIDs, p.ChatID)
	}

	participantsByChat, err := uc.store.BatchParticipantsForChats(ctx, chatIDs)
	if err != nil {
		uc.log.Error("Failed to batch load participants",
			withFields("error", err.Error())...)
		return emptyWithErr(fmt.Errorf("%s: %w", op, err))
	}

	for i := range previews {
		previews[i].Participants = participantsByChat[previews[i].ChatID]
	}

	previewsResp := make([]dto.ChatPreview, 0, len(previews))
	for _, preview := range previews {
		chatPreviewResp := dto.ChatPreview{
			ChatID:      preview.ChatID.String(),
			Name:        preview.Name,
			UnreadCount: preview.UnreadCount,
			AvatarURL:   preview.AvatarURL,
		}

		participantsResp := make([]dto.Participant, 0, len(preview.Participants))
		for _, participant := range preview.Participants {
			participantsResp = append(participantsResp, dto.Participant{
				UserID:    participant.UserID.String(),
				Username:  participant.Username,
				AvatarURL: participant.AvatarURL,
				IsOnline:  participant.IsOnline,
			})
		}
		chatPreviewResp.Participants = participantsResp

		if preview.LastMessage.Text != "" {
			chatPreviewResp.LastMessage = dto.Message{
				Content:   preview.LastMessage.Text,
				UserID:    preview.LastMessage.SenderID.String(),
				Timestamp: preview.LastMessage.Timestamp.Format(time.RFC3339),
			}
		}

		previewsResp = append(previewsResp, chatPreviewResp)
	}

	uc.log.Info("Successfully got user chats", withFields("count", len(previewsResp))...)
	return dto.GetUserChatPreviewResponse{ChatPreviews: previewsResp}, nil
}
func emptyWithErr(err error) (dto.GetUserChatPreviewResponse, error) {
	return dto.GetUserChatPreviewResponse{}, err
}
