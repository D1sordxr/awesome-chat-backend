package process

import (
	"awesome-chat/internal/domain/core/shared/broker/entity"
	"context"
	"errors"
	"fmt"

	"awesome-chat/internal/domain/core/shared/outbox/filters"
	"awesome-chat/internal/domain/core/shared/outbox/ports"
	"awesome-chat/internal/domain/core/shared/outbox/vo"
	sharedPorts "awesome-chat/internal/domain/core/shared/ports"
)

type UseCase struct {
	txManager sharedPorts.TransactionManager
	store     ports.ProcessStore
	producer  sharedPorts.Producer
}

func NewUseCase(
	store ports.ProcessStore,
	txManager sharedPorts.TransactionManager,
	producer sharedPorts.Producer,
) *UseCase {
	return &UseCase{
		store:     store,
		txManager: txManager,
		producer:  producer,
	}
}

func (uc *UseCase) Process(ctx context.Context, filter filters.GetOutbox) error {
	ctx, err := uc.txManager.BeginAndInjectTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx error: %w", err)
	}
	defer func() { _ = uc.txManager.RollbackTx(ctx) }()

	messages, err := uc.store.GetMessagesForUpdate(ctx, filter)
	if err != nil {
		return fmt.Errorf("get message error: %w", err)
	}

	var errs []error
	for _, msg := range messages {
		if err = uc.producer.Publish(ctx, entity.Message{
			Key:   []byte(msg.OutboxID.String()),
			Value: msg.Payload,
		}); err != nil {
			if updateErr := uc.store.UpdateStatus(ctx, filters.SetOutboxStatus{
				OutboxID: msg.OutboxID,
				Status:   vo.StatusFailed,
			}); updateErr != nil {
				errs = append(errs, fmt.Errorf("mark failed (id=%s): %w", msg.OutboxID, updateErr))
			}
			errs = append(errs, fmt.Errorf("send (id=%s): %w", msg.OutboxID, err))
			continue
		}

		if err = uc.store.UpdateStatus(ctx, filters.SetOutboxStatus{
			OutboxID: msg.OutboxID,
			Status:   vo.StatusProcessed,
		}); err != nil {
			errs = append(errs, fmt.Errorf("mark processed (id=%s): %w", msg.OutboxID, err))
		}
	}

	if err = uc.txManager.CommitTx(ctx); err != nil {
		if len(errs) > 0 {
			err = fmt.Errorf("commit tx error: %w: processed with errors: %w", err, errors.Join(errs...))
		}
		return fmt.Errorf("commit tx error: %w", err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("processed with errors: %w", errors.Join(errs...))
	}

	return nil
}
