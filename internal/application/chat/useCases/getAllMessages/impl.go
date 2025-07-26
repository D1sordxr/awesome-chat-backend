package getAllMessages

import (
	"awesome-chat/internal/application/chat/dto"
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/chat/ports"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type ChatGetAllMessagesUseCase struct {
	log   appPorts.Logger
	store ports.GetAllMessagesStore
}

func NewChatGetAllMessagesUseCase(
	log appPorts.Logger,
	store ports.GetAllMessagesStore,
) *ChatGetAllMessagesUseCase {
	return &ChatGetAllMessagesUseCase{
		log:   log,
		store: store,
	}
}

func (uc *ChatGetAllMessagesUseCase) Execute(ctx context.Context, chatID dto.ChatID) (dto.AllMessages, error) {
	const op = "ChatGetAllMessagesUseCase.SetupChatPreviews"
	withFields := func(args ...any) []any {
		return append([]any{"op", op, "chatID", string(chatID)}, args...)
	}

	uc.log.Info("Attempting to get all chat messages", withFields()...)

	id, err := uuid.Parse(string(chatID))
	if err != nil {
		uc.log.Error("Failed to parse chat ID", withFields("error", err.Error())...)
		return dto.AllMessages{}, fmt.Errorf("%s: %w", op, err)
	}

	messages, err := uc.store.Execute(ctx, id)
	if err != nil {
		uc.log.Error("Failed to get all chat messages", withFields("error", err.Error())...)
		return dto.AllMessages{}, fmt.Errorf("%s: %w", op, err)
	}

	msgResp := make(dto.AllMessages, 0, len(messages))
	for _, msg := range messages {
		msgResp = append(msgResp, dto.Message{
			Text:      msg.Text,
			SenderID:  msg.SenderID.String(),
			Timestamp: msg.Timestamp.String(),
		})
	}

	uc.log.Info("Successfully got all chat messages", withFields()...)
	return msgResp, nil
}
