package stream

import (
	appPorts "awesome-chat/internal/domain/app/ports"
	"awesome-chat/internal/domain/core/message/ports/worker"
	conn "awesome-chat/internal/infrastructure/redis"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type SubscriberImpl struct {
	log appPorts.Logger

	client     *conn.Connection
	streamName string
	groupName  string
	consumerID string

	streamPipe worker.MessagePipe[redis.XMessage]
	done       chan struct{}
}

func NewSubscriberImpl(
	log appPorts.Logger,
	conn *conn.Connection,
	streamPipe worker.MessagePipe[redis.XMessage],
	streamName string,
	groupName string,
	consumerID string,
) *SubscriberImpl {
	return &SubscriberImpl{
		log:        log,
		client:     conn,
		streamName: streamName,
		groupName:  groupName,
		consumerID: consumerID,
		streamPipe: streamPipe,
		done:       make(chan struct{}),
	}
}

// Ack sets message as processed
func (s *SubscriberImpl) Ack(ctx context.Context, msgID string) error {
	return s.client.XAck(ctx, s.streamName, s.groupName, msgID).Err()
}

func (s *SubscriberImpl) readGroup(ctx context.Context) error {
	messages, err := s.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    s.groupName,
		Consumer: s.consumerID,
		Streams:  []string{s.streamName, ">"},
		Count:    10,
		Block:    5 * time.Second,
		NoAck:    false,
	}).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}

	if len(messages) == 0 || len(messages[0].Messages) == 0 {
		return nil
	}

	for _, msg := range messages[0].Messages {
		select {
		case s.streamPipe.GetWriteChan() <- msg:
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}

func (s *SubscriberImpl) groupCheck(ctx context.Context) error {
	err := s.client.XGroupCreateMkStream(ctx, s.streamName, s.groupName, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {

		return err
	}
	s.log.Warn("Stream already exists",
		"stream_name", s.streamName, "group_name", s.groupName)
	return nil
}

func (s *SubscriberImpl) Start(ctx context.Context) error {
	if err := s.groupCheck(ctx); err != nil {
		s.log.Error("Group check failed", "error", err.Error())
		return err
	}

	s.log.Info(`Group check successful. Starting subscriber...`,
		"streamName", s.streamName,
		"groupName", s.groupName,
		"consumerID", s.consumerID,
	)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-s.done:
			return nil
		default:
			if err := s.readGroup(ctx); err != nil {
				if !errors.Is(err, redis.Nil) {
					time.Sleep(time.Second)
					continue
				}
				s.log.Error("Subscriber read error", "error", err.Error())
			}
		}
	}
}

func (s *SubscriberImpl) Shutdown(_ context.Context) error {
	close(s.done)
	s.streamPipe.Close()
	return nil
}
