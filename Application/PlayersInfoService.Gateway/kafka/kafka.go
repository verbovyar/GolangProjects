package kafka

import (
	"github.com/Shopify/sarama"
	"log"
)

var brokers = []string{"127.0.0.1:9092"} //TODO move in config file

func NewProducer() (sarama.SyncProducer, error) {
	conf := sarama.NewConfig()
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.RequiredAcks = sarama.WaitForLocal
	conf.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, conf)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func NewConsumerGroup() sarama.ConsumerGroup {
	conf := sarama.NewConfig()
	conf.Consumer.Return.Errors = true
	conf.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, "startConsuming", conf)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return consumerGroup
}
