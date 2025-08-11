package broadcast

import (
	"awesome-chat/internal/application/message/dto"
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/message/vo"
	"awesome-chat/internal/domain/core/shared/ports"
	"awesome-chat/internal/domain/core/shared/ports/ws"
	"awesome-chat/internal/infrastructure/ws/chathub"
	"context"
	"fmt"
	"time"
)

type MessageBroadcastWithPubImpl struct {
	log appPorts.Logger
	pub ports.StreamPublisher
	br  ws.MessageBroadcaster
}

func NewMessageBroadcastWithPubImpl(
	log appPorts.Logger,
	pub ports.StreamPublisher,
	br ws.MessageBroadcaster,
) *MessageBroadcastWithPubImpl {
	return &MessageBroadcastWithPubImpl{
		log: log,
		pub: pub,
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
		}, fields...)
	}
	m.log.Info("Starting operation", withFields()...)

	timestamp := time.Now()

	if err := m.pub.Publish(ctx, vo.StreamMessage{
		Event:     vo.SentMessageEvent,
		UserID:    req.UserID,
		ChatID:    req.ChatID,
		Content:   req.Content,
		Timestamp: timestamp,
	}.ToMap()); err != nil {
		m.log.Error("Failed to publish message.",
			withFields("error", err.Error())...)
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := m.br.Broadcast(ctx, chathub.Message{ // TODO: prefer to replace into domain
		UserID:    req.UserID,
		ChatID:    req.ChatID,
		Content:   req.Content,
		Timestamp: timestamp.String(),
	}); err != nil {
		m.log.Error("Failed to broadcast message. Message will be saved to DB.",
			withFields("error", err.Error())...)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
