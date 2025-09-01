package ioc

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gotomicro/ego/core/econf"
)

func InitKafkaConsumer(groupID string) *kafka.Consumer {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := econf.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.Addr,
		"group.id":           groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "false",
	})
	if err != nil {
		panic(err)
	}
	return consumer
}
