package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"awesome-chat/internal/application/message/command"
	"awesome-chat/internal/application/message/port/mocks"
	chatErrors "awesome-chat/internal/domain/core/chat/errors"
	"awesome-chat/internal/domain/core/message/entity"
	outboxEntity "awesome-chat/internal/domain/core/shared/outbox/entity"
	outboxVO "awesome-chat/internal/domain/core/shared/outbox/vo"
)

func TestSendDeniesNonMember(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	chats := mocks.NewMockChatValidator(ctrl)
	repo := mocks.NewMockRepository(ctrl)

	senderID, chatID := uuid.New(), uuid.New()

	chats.EXPECT().IsMember(gomock.Any(), chatID, senderID).Return(false, nil)

	send := command.NewSend(repo, mocks.NewMockOutbox(ctrl), chats, mocks.NewMockTxManager(ctrl))

	_, err := send.Handle(context.Background(), command.SendIn{
		SenderID: senderID,
		ChatID:   chatID,
		Content:  "hello",
	})

	if !errors.Is(err, chatErrors.ErrChatAccessDenied) {
		t.Fatalf("err = %v, want %v", err, chatErrors.ErrChatAccessDenied)
	}
}

func TestSendWritesMessageAndOutboxInOneTransaction(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	chats := mocks.NewMockChatValidator(ctrl)
	repo := mocks.NewMockRepository(ctrl)
	outbox := mocks.NewMockOutbox(ctrl)
	tx := mocks.NewMockTxManager(ctrl)

	senderID, chatID := uuid.New(), uuid.New()
	saved := entity.Message{ID: 42, UserID: senderID, ChatID: chatID, Content: "hello"}

	chats.EXPECT().IsMember(gomock.Any(), chatID, senderID).Return(true, nil)

	tx.EXPECT().
		WithTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	repo.EXPECT().
		Save(gomock.Any(), entity.Message{UserID: senderID, ChatID: chatID, Content: "hello"}).
		Return(saved, nil)

	outbox.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, record outboxEntity.Outbox) error {
			if record.EntityName != outboxVO.MessageEntity {
				t.Errorf("entity name = %q, want %q", record.EntityName, outboxVO.MessageEntity)
			}
			if record.Status != outboxVO.StatusPending {
				t.Errorf("status = %q, want %q", record.Status, outboxVO.StatusPending)
			}
			if len(record.Payload) == 0 {
				t.Error("outbox payload is empty")
			}

			return nil
		})

	send := command.NewSend(repo, outbox, chats, tx)

	out, err := send.Handle(context.Background(), command.SendIn{
		SenderID: senderID,
		ChatID:   chatID,
		Content:  "hello",
	})
	if err != nil {
		t.Fatalf("send: %v", err)
	}

	if out.Message.ID != saved.ID {
		t.Errorf("message id = %d, want %d", out.Message.ID, saved.ID)
	}
}
