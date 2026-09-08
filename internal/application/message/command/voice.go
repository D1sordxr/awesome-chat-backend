package command

import (
	"time"

	"context"

	"github.com/google/uuid"

	"awesome-chat/internal/application/message/port"
	chatErrors "awesome-chat/internal/domain/core/chat/errors"
	"awesome-chat/internal/domain/core/message/entity"
	"awesome-chat/internal/domain/core/message/vo"
)

type CreateVoiceUpload struct {
	uploads port.UploadService
	chats   port.ChatValidator
}

func NewCreateVoiceUpload(uploads port.UploadService, chats port.ChatValidator) *CreateVoiceUpload {
	return &CreateVoiceUpload{uploads: uploads, chats: chats}
}

type CreateVoiceUploadIn struct {
	CallerID uuid.UUID
	ChatID   uuid.UUID
}

type CreateVoiceUploadOut struct {
	ObjectKey string
	UploadURL string
	ExpiresAt time.Time
}

func (c *CreateVoiceUpload) Handle(ctx context.Context, in CreateVoiceUploadIn) (CreateVoiceUploadOut, error) {
	isMember, err := c.chats.IsMember(ctx, in.ChatID, in.CallerID)
	if err != nil {
		return CreateVoiceUploadOut{}, err
	}

	if !isMember {
		return CreateVoiceUploadOut{}, chatErrors.ErrChatAccessDenied
	}

	objectKey := in.ChatID.String() + "/" + uuid.NewString()

	uploadURL, expiresAt, err := c.uploads.PresignedPut(ctx, objectKey)
	if err != nil {
		return CreateVoiceUploadOut{}, err
	}

	return CreateVoiceUploadOut{
		ObjectKey: objectKey,
		UploadURL: uploadURL,
		ExpiresAt: expiresAt,
	}, nil
}

type SendVoice struct {
	repo  port.VoiceRepository
	chats port.ChatValidator
	tx    port.TxManager
}

func NewSendVoice(repo port.VoiceRepository, chats port.ChatValidator, tx port.TxManager) *SendVoice {
	return &SendVoice{repo: repo, chats: chats, tx: tx}
}

type SendVoiceIn struct {
	SenderID        uuid.UUID
	ChatID          uuid.UUID
	ObjectKey       string
	DurationSeconds int
}

type SendVoiceOut struct {
	Message entity.Message
	Voice   vo.SaveVoiceData
}

func (c *SendVoice) Handle(ctx context.Context, in SendVoiceIn) (SendVoiceOut, error) {
	isMember, err := c.chats.IsMember(ctx, in.ChatID, in.SenderID)
	if err != nil {
		return SendVoiceOut{}, err
	}

	if !isMember {
		return SendVoiceOut{}, chatErrors.ErrChatAccessDenied
	}

	voice := vo.SaveVoiceData{
		UserID:          in.SenderID,
		ChatID:          in.ChatID,
		ObjectKey:       in.ObjectKey,
		DurationSeconds: in.DurationSeconds,
	}

	var saved entity.Message

	if err = c.tx.WithTransaction(ctx, func(txCtx context.Context) error {
		saved, err = c.repo.SaveVoice(txCtx, voice)

		return err
	}); err != nil {
		return SendVoiceOut{}, err
	}

	return SendVoiceOut{Message: saved, Voice: voice}, nil
}
