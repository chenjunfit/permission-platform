package mqx

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

type Consumer interface {
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Assignment() (partitions []kafka.TopicPartition, err error)
	Pause(partitions []kafka.TopicPartition) (err error)
	Resume(partitions []kafka.TopicPartition) (err error)
	Poll(timeoutMs int) (event kafka.Event)
	CommitMessage(m *kafka.Message) ([]kafka.TopicPartition, error)
}
