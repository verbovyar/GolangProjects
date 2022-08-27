package kafka

import "github.com/Shopify/sarama"

var brokers = []string{"127.0.0.1:9092"} //TODO move in config file

func NewProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	producer, err := sarama.NewSyncProducer(brokers, config)

	return producer, err
}

// TODO create consumer
