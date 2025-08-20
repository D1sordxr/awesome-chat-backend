package sendVoice

import (
	"awesome-chat/internal/application/message/dto"
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/message/ports/store"
	"awesome-chat/internal/domain/core/message/vo"
	"awesome-chat/internal/domain/core/shared/ports"
	"awesome-chat/internal/domain/core/shared/ports/s3"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type MessageSendVoiceUseCase struct {
	log       appPorts.Logger
	txManager ports.TransactionManager
	store     store.SaveVoice
	s3Storage s3.Storage
	s3UrlSvc  s3.URLService
}

func NewMessageSendVoiceUseCase(
	log appPorts.Logger,
	txManager ports.TransactionManager,
	store store.SaveVoice,
	s3Storage s3.Storage,
	s3UrlSvc s3.URLService,
) *MessageSendVoiceUseCase {
	return &MessageSendVoiceUseCase{
		log:       log,
		txManager: txManager,
		store:     store,
		s3Storage: s3Storage,
		s3UrlSvc:  s3UrlSvc,
	}
}

func (uc *MessageSendVoiceUseCase) Execute(
	ctx context.Context,
	req *dto.SendVoiceRequest,
) (
	dto.SendVoiceResponse,
	error,
) {
	const op = "MessageSendVoiceUseCase.Execute"
	withFields := func(fields ...any) []any {
		return append([]any{"operation", op, "chatID", req.ChatID}, fields...)
	}

	uc.log.Info("Attempting to send voice message", withFields()...)

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return dto.SendVoiceResponse{}, fmt.Errorf("%s, %w", op, err)
	}
	chatID, err := uuid.Parse(req.ChatID)
	if err != nil {
		return dto.SendVoiceResponse{}, fmt.Errorf("%s, %w", op, err)
	}
	voiceID := uuid.New()

	if err = uc.s3Storage.Add(ctx, voiceID.String(), []byte(req.Blob)); err != nil {
		return dto.SendVoiceResponse{}, fmt.Errorf("%w: %s", err, op)
	}
	defer func() {
		if err != nil {
			_ = uc.s3Storage.Delete(ctx, voiceID.String())
		}
	}()

	audioURL, err := uc.s3UrlSvc.GenerateURL(ctx, voiceID.String())
	if err != nil {
		return dto.SendVoiceResponse{}, fmt.Errorf("%w: %s", err, op)
	}

	// TODO: generate waveform

	ctx, err = uc.txManager.BeginAndInjectTx(ctx)
	if err != nil {
		return dto.SendVoiceResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			_ = uc.txManager.RollbackTx(ctx)
		}
	}()

	if err = uc.store.Execute(ctx, vo.SaveVoiceData{
		UserID:   userID,
		ChatID:   chatID,
		AudioURL: audioURL,
		Duration: req.Duration,
		Waveform: nil, // TODO generate
	}); err != nil {
		return dto.SendVoiceResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	if err = uc.txManager.CommitTx(ctx); err != nil {
		return dto.SendVoiceResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	uc.log.Info("Successfully sent voice message", withFields()...)

	return dto.SendVoiceResponse{}, nil
}
