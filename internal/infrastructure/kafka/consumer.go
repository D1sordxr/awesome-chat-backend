package kafka

import (
	"awesome-chat/internal/domain/core/shared/broker/entity"
	cfg "awesome-chat/internal/infrastructure/config/kafka"
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	r *kafka.Reader
}

func NewConsumer(cfg *cfg.Config) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		GroupID: "",
		// MinBytes: 10e3, // 10KB
		// MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		r: r,
	}
}

func (c *Consumer) Receive(ctx context.Context) (entity.Message, error) {
	msg, err := c.r.ReadMessage(ctx)
	if err != nil {
		return entity.Message{}, err
	}
	return entity.Message{
		Key:   msg.Key,
		Value: msg.Value,
	}, nil
}
