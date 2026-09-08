package message

import (
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/shared/ports"
	"context"
	"sync"
)

const (
	workers         = 10
	messageChanSize = 100
)

type useCase interface {
	Execute(ctx context.Context, payload []byte) error
}

type FastSaveAndBroadcastMessageHandler struct {
	log      appPorts.Logger
	sub      ports.Subscriber
	uc       useCase
	msgChan  chan []byte
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewFastSaveAndBroadcastMessageHandler(
	log appPorts.Logger,
	sub ports.Subscriber,
	uc useCase,
) *FastSaveAndBroadcastMessageHandler {
	return &FastSaveAndBroadcastMessageHandler{
		log:      log,
		sub:      sub,
		uc:       uc,
		msgChan:  make(chan []byte, messageChanSize),
		stopChan: make(chan struct{}),
	}
}

func (h *FastSaveAndBroadcastMessageHandler) Start(ctx context.Context) error {
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch, err := h.sub.GetSubChannel(subCtx)
	if err != nil {
		return err
	}

	h.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go h.worker(subCtx)
	}

	for {
		select {
		case <-ctx.Done():
			close(h.stopChan)
			h.wg.Wait()
			return nil
		case msg, ok := <-ch:
			if !ok {
				close(h.stopChan)
				h.wg.Wait()
				return nil
			}
			select {
			case h.msgChan <- msg:
			case <-ctx.Done():
				close(h.stopChan)
				h.wg.Wait()
				return nil
			}
		}
	}
}

func (h *FastSaveAndBroadcastMessageHandler) worker(ctx context.Context) {
	defer h.wg.Done()

	for {
		select {
		case <-h.stopChan:
			return
		case msg, ok := <-h.msgChan:
			if !ok {
				return
			}
			h.processMessage(ctx, msg)
		}
	}
}

func (h *FastSaveAndBroadcastMessageHandler) processMessage(ctx context.Context, msg []byte) {
	if err := h.uc.Execute(ctx, msg); err != nil {
		h.log.Error("message save and broadcast error", "error", err.Error())
		return
	}
}
