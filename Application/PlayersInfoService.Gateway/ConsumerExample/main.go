package main

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

type Consumer struct {
}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			log.Print("Done")
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Print("Data channel closed")
				return nil
			}
			log.Printf("partition: %v, data: %v", msg.Partition, string(msg.Value))
			session.MarkMessage(msg, "")
		}
	}
}

func main() {
	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(brokers, "startConsuming", config)
	if err != nil {
		log.Fatalf(err.Error())
	}
	ctx := context.Background()
	consumer := &Consumer{}
	for {
		if err = client.Consume(ctx, []string{"AddRequest"}, consumer); err != nil {
			log.Printf("on consume: %v", err)
			time.Sleep(time.Second * 10)
		}
	}
}
