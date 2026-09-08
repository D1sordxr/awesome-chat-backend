package command

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"awesome-chat/internal/application/message/port"
	chatErrors "awesome-chat/internal/domain/core/chat/errors"
	"awesome-chat/internal/domain/core/message/entity"
	outboxEntity "awesome-chat/internal/domain/core/shared/outbox/entity"
	outboxVO "awesome-chat/internal/domain/core/shared/outbox/vo"
)

type Send struct {
	repo   port.Repository
	outbox port.Outbox
	chats  port.ChatValidator
	tx     port.TxManager
}

func NewSend(
	repo port.Repository,
	outbox port.Outbox,
	chats port.ChatValidator,
	tx port.TxManager,
) *Send {
	return &Send{repo: repo, outbox: outbox, chats: chats, tx: tx}
}

type SendIn struct {
	SenderID uuid.UUID
	ChatID   uuid.UUID
	Content  string
}

type SendOut struct {
	Message entity.Message
}

func (c *Send) Handle(ctx context.Context, in SendIn) (SendOut, error) {
	isMember, err := c.chats.IsMember(ctx, in.ChatID, in.SenderID)
	if err != nil {
		return SendOut{}, err
	}

	if !isMember {
		return SendOut{}, chatErrors.ErrChatAccessDenied
	}

	var saved entity.Message

	if err = c.tx.WithTransaction(ctx, func(txCtx context.Context) error {
		if saved, err = c.repo.Save(txCtx, entity.Message{
			UserID:  in.SenderID,
			ChatID:  in.ChatID,
			Content: in.Content,
		}); err != nil {
			return err
		}

		payload, marshalErr := json.Marshal(saved)
		if marshalErr != nil {
			return fmt.Errorf("marshal outbox payload: %w", marshalErr)
		}

		return c.outbox.Save(txCtx, outboxEntity.Outbox{
			OutboxID:   uuid.New(),
			EntityName: outboxVO.MessageEntity,
			Status:     outboxVO.StatusPending,
			Payload:    payload,
		})
	}); err != nil {
		return SendOut{}, err
	}

	return SendOut{Message: saved}, nil
}
