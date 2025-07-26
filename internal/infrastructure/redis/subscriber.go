package redis

import (
	cfg "awesome-chat/internal/infrastructure/config/redis"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

type Subscriber struct {
	channel string
	client  *redis.Client
}

func NewSubscriber(cfg *cfg.Config) *Subscriber {
	return &Subscriber{
		channel: cfg.GetChannel(),
		client: redis.NewClient(&redis.Options{
			Addr: cfg.GetClientAddress(),
		}),
	}
}

func (s *Subscriber) GetSubChannel(ctx context.Context) (<-chan []byte, error) {
	pubSub := s.client.Subscribe(ctx, s.channel)

	if _, err := pubSub.Receive(ctx); err != nil {
		_ = pubSub.Close()
		return nil, errors.New("failed to subscribe: " + err.Error())
	}

	msgChan := make(chan []byte, 100)

	go func() {
		defer func() {
			close(msgChan)
			_ = pubSub.Close()
		}()

		redisChan := pubSub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-redisChan:
				if !ok {
					return
				}
				select {
				case msgChan <- []byte(msg.Payload):
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return msgChan, nil
}
