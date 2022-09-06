package kafka

import (
	"github.com/Shopify/sarama"
	"log"
)

var brokers = []string{"127.0.0.1:9092"}

func NewProducer() sarama.SyncProducer {
	conf := sarama.NewConfig()
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.RequiredAcks = sarama.WaitForLocal
	conf.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, conf)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return producer
}

func NewConsumerGroup() sarama.ConsumerGroup {
	conf := sarama.NewConfig()
	conf.Consumer.Return.Errors = true
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, "startConsuming", conf)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return consumerGroup
}
