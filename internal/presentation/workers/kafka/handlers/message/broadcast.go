package message

import (
	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/shared/broker/entity"
	sharedPorts "awesome-chat/internal/domain/core/shared/ports"
	"context"
	"time"
)

const (
	handler       = "BroadcastMessage"
	tickerTime    = 100 * time.Millisecond
	broadcastTime = 5 * time.Second
)

type useCase interface {
	Broadcast(ctx context.Context, message entity.Message) error
}

type BroadcastMessage struct {
	log      ports.Logger
	uc       useCase
	consumer sharedPorts.Consumer
	ticker   *time.Ticker
}

func NewBroadcastMessage(
	log ports.Logger,
	consumer sharedPorts.Consumer,
	uc useCase,
) *BroadcastMessage {
	return &BroadcastMessage{
		log:      log,
		uc:       uc,
		consumer: consumer,
		ticker:   time.NewTicker(tickerTime),
	}
}

func (b *BroadcastMessage) Start(ctx context.Context) error {
	defer b.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			b.log.Info("handler stopped via context")
			return nil
		case <-b.ticker.C:
			if err := b.receiveAndBroadcast(ctx); err != nil {
				b.log.Error(err.Error(), "handler", handler)
			}
		}
	}
}

func (b *BroadcastMessage) receiveAndBroadcast(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, broadcastTime)
	defer cancel()

	msg, err := b.consumer.Receive(ctx)
	if err != nil {
		return err
	}

	if err = b.uc.Broadcast(ctx, msg); err != nil {
		return err
	}

	return nil
}
