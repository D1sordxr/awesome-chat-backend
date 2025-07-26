package send

import (
	"awesome-chat/internal/application/message/dto"
	"awesome-chat/internal/domain/core/message/ports"
	"bytes"
	"context"
	"errors"
	"net/http"
	"time"
)

type MessageSendSyncUseCase struct {
	repo           ports.Repository
	entityCreator  ports.EntityCreator
	wsServerClient *http.Client
	wsServerURL    string
}

func NewMessageSendSyncUseCase(
	repo ports.Repository,
	entityCreator ports.EntityCreator,
	wsServerURL string,
) *MessageSendSyncUseCase {
	return &MessageSendSyncUseCase{
		repo:           repo,
		entityCreator:  entityCreator,
		wsServerClient: http.DefaultClient,
		wsServerURL:    wsServerURL,
	}
}

func (uc *MessageSendSyncUseCase) Execute(
	ctx context.Context,
	req dto.SendSyncRequest,
	raw []byte,
) error {
	entity := uc.entityCreator.Do(
		req.UserID,
		req.ChatID,
		req.Content,
	)

	if err := uc.repo.Save(ctx, entity); err != nil {
		return err
	}

	reqCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	msgBuff := bytes.NewBuffer(raw)

	r, err := http.NewRequestWithContext(reqCtx, http.MethodPost, uc.wsServerURL, msgBuff)
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")

	resp, err := uc.wsServerClient.Do(r)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return errors.New(resp.Status)
	}

	return nil
}
