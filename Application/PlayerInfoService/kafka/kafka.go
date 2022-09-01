package kafka

import (
	"github.com/Shopify/sarama"
)

func NewProducer() (sarama.SyncProducer, error) {
	var brokers = []string{"127.0.0.1:9092"}

	conf := sarama.NewConfig()
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.RequiredAcks = sarama.WaitForLocal
	conf.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, conf)

	return producer, err
}
