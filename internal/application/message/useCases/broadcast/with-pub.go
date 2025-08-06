package broadcast

import (
	"awesome-chat/internal/application/message/dto"
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/shared/ports/ws"
	"awesome-chat/internal/infrastructure/ws/chathub"
	"context"
	"fmt"
)

type MessageBroadcastWithPubImpl struct {
	log appPorts.Logger
	// redis client
	br ws.MessageBroadcaster
}

func NewMessageBroadcastWithPubImpl(
	log appPorts.Logger,
	br ws.MessageBroadcaster,
) *MessageBroadcastWithPubImpl {
	return &MessageBroadcastWithPubImpl{
		log: log,
		br:  br,
	}
}

func (m *MessageBroadcastWithPubImpl) Execute(
	ctx context.Context,
	req dto.BroadcastWithPubRequest,
) error {
	const op = "MessageBroadcastWithPub.Execute"
	withFields := func(fields ...any) []any {
		return append([]any{
			"operation", op,
			"user_id", req.UserID,
			"chat_id", req.ChatID,
			"content_length", len(req.Content),
		}, fields...)
	}
	m.log.Info("Starting operation", withFields()...)

	// TODO: redis client
	m.log.Warn("Redis client not implemented. Message won't be saved", withFields()...)

	msg := chathub.Message{
		UserID:  req.UserID,
		ChatID:  req.ChatID,
		Content: req.Content,
	}
	if err := m.br.Broadcast(ctx, msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
