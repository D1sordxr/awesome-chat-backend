package getForChatWithFilter

import (
	"awesome-chat/internal/application/message/dto"
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/message/ports"
	"awesome-chat/internal/domain/core/message/vo"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type MessageGetForChatWithFilterUseCase struct {
	log   appPorts.Logger
	store ports.GetForChatWithFilterStore
}

func NewMessageGetForChatWithFilterUseCase(
	log appPorts.Logger,
	store ports.GetForChatWithFilterStore,
) *MessageGetForChatWithFilterUseCase {
	return &MessageGetForChatWithFilterUseCase{
		log:   log,
		store: store,
	}
}

func (uc *MessageGetForChatWithFilterUseCase) Execute(
	ctx context.Context,
	req dto.GetForChatWithFilterRequest,
) (
	dto.GetForChatWithFilterResponse,
	error,
) {
	const op = "MessageGetForChatWithFilterUseCase.Execute"
	withFields := func(args ...any) []any {
		return append([]any{"op", op, "chat_id", req.ChatID}, args...)
	}

	uc.log.Info("Attempting to get messages with filter", withFields()...)

	var (
		chatID uuid.UUID
		limit  int
		// offset int
		cursor int
	)

	chatID, err := uuid.Parse(req.ChatID)
	if err != nil {
		return dto.GetForChatWithFilterResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	if req.Limit <= 0 {
		limit = 100
	}
	cursor = req.Cursor

	filter := vo.ReadFilter{
		ChatID: chatID,
		Limit:  limit,
		Cursor: cursor,
		// Offset: 0,
	}

	messages, err := uc.store.Execute(ctx, filter)
	if err != nil {
		return dto.GetForChatWithFilterResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	filteredMessage := make([]dto.FilteredMessage, 0, len(messages))
	for _, message := range messages {
		msg := dto.FilteredMessage{
			ID:        message.ID,
			Text:      message.Text,
			SenderID:  message.SenderID.String(),
			Timestamp: message.Timestamp.String(),
		}
		filteredMessage = append(filteredMessage, msg)
	}

	lastCursor := 0
	if len(filteredMessage) > 0 {
		lastCursor = filteredMessage[len(filteredMessage)-1].ID
	}

	if len(messages) == 0 {
		uc.log.Warn("No messages found", withFields("filter", filter)...)
	}

	uc.log.Info("Successfully got messages with filter", withFields("filter", filter)...)

	return dto.GetForChatWithFilterResponse{
		AllMessages: filteredMessage,
		Count:       len(messages),
		LastCursor:  lastCursor,
	}, nil
}
