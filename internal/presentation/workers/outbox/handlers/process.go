package handlers

import (
	"awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/shared/outbox/filters"
	"context"
	"time"
)

const tickerTime = 1 * time.Second

type useCase interface {
	Process(ctx context.Context, getFilter filters.GetOutbox) error
}

type ProcessHandler struct {
	uc     useCase
	log    ports.Logger
	filter filters.GetOutbox
	ticker *time.Ticker
}

func NewHandler(
	log ports.Logger,
	uc useCase,
	filter filters.GetOutbox,
) *ProcessHandler {
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	return &ProcessHandler{
		log:    log,
		uc:     uc,
		filter: filter,
		ticker: time.NewTicker(tickerTime),
	}
}

func (p *ProcessHandler) Start(ctx context.Context) error {
	defer p.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.log.Info("handler stopped via context")
			return nil
		case <-p.ticker.C:
			p.processOutbox(ctx)
		}
	}
}

func (p *ProcessHandler) processOutbox(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := p.uc.Process(ctx, p.filter); err != nil {
		p.log.Error("outbox processing failed",
			"error", err.Error(),
			"entity", p.filter.EntityName,
		)
	} else {
		p.log.Debug("outbox batch processed",
			"entity", p.filter.EntityName,
		)
	}
}
