package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

var brokers = []string{"127.0.0.1:9092"} //TODO move in config file

func NewProducer() (sarama.SyncProducer, error) {
	conf := sarama.NewConfig()
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.RequiredAcks = sarama.WaitForLocal
	conf.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, conf)

	return producer, err
}

func NewConsumer() sarama.Consumer {
	conf := sarama.NewConfig()
	conf.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, conf)
	if err != nil {
		fmt.Println(err)
	}

	return consumer
}
