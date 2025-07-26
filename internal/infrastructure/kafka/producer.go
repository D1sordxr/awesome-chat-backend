package kafka

import (
	"awesome-chat/internal/domain/core/shared/broker/entity"
	cfg "awesome-chat/internal/infrastructure/config/kafka"
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Topic  string
	Writer *kafka.Writer
}

func NewProducer(config *cfg.Config) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(config.Brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		Topic:  config.Topic,
		Writer: writer,
	}
}

func (p *Producer) Publish(ctx context.Context, message entity.Message) error {
	return p.Writer.WriteMessages(ctx, kafka.Message{
		Topic: p.Topic,
		Key:   message.Key,
		Value: message.Value,
	})
}

func (p *Producer) Shutdown(_ context.Context) error {
	return p.Writer.Close()
}
