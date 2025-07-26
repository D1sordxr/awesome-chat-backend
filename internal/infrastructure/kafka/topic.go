package kafka

import (
	"github.com/segmentio/kafka-go"

	config "awesome-chat/internal/infrastructure/config/kafka"
)

const (
	partitions        = 3
	replicationFactor = 1
)

func CreateTopic(cfg *config.Config) {
	conn, err := kafka.Dial("tcp", cfg.Brokers[0])
	if err != nil {
		panic("Failed to connect to Kafka" + err.Error())
	}
	defer func() { _ = conn.Close() }()

	topicConfigs := kafka.TopicConfig{
		Topic:             cfg.Topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}
	if err = conn.CreateTopics(topicConfigs); err != nil {
		panic("Failed to create topics:" + err.Error())
	}
}
