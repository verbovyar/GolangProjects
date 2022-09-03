package addHandler

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"modules/internal/handlers"
	pbGoFiles2 "modules/internal/infrastructure/playersInfoServiceClient/api/pbGoFiles"
	"time"
)

type AddHandler struct {
	producer      sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup
}

func NewAddHandler(producer sarama.SyncProducer, consumerGroup sarama.ConsumerGroup) *AddHandler {
	return &AddHandler{
		producer:      producer,
		consumerGroup: consumerGroup,
	}
}

func (c *AddHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *AddHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *AddHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var response pbGoFiles2.AddResponse
		err := json.Unmarshal(msg.Value, &response)
		if err != nil {
			log.Printf("income data %v: %v", string(msg.Value), err)
			continue
		}

		handlers.AddResponse{Id: response.Id}
	}

	return nil
}

func AddClaim(addHandler *AddHandler) {
	for {
		if err := addHandler.consumerGroup.Consume(context.Background(), []string{"AddResponse"}, addHandler); err != nil {
			log.Printf("on consume: %v", err)
			time.Sleep(time.Second * 10)
		}
	}
}
