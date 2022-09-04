package deleteHandler

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/go-telegram-bot-api/telegram-bot-api/api/apiPb"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/kafkaHandlers/structs"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/repositories/interfaces"
	"log"
	"time"
)

type DeleteHandler struct {
	producer      sarama.SyncProducer
	repository    interfaces.Repository
	consumerGroup sarama.ConsumerGroup
}

func NewDeleteHandler(producer sarama.SyncProducer, repository interfaces.Repository, consumerGroup sarama.ConsumerGroup) *DeleteHandler {
	return &DeleteHandler{
		producer:      producer,
		repository:    repository,
		consumerGroup: consumerGroup,
	}
}

func (c *DeleteHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *DeleteHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *DeleteHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var deleteRequest structs.DeleteRequest
		err := json.Unmarshal(msg.Value, &deleteRequest)
		if err != nil {
			log.Printf("income data %v: %v", string(msg.Value), err)
			continue
		}

		err = c.repository.Delete(uint(deleteRequest.Id))
		if err != nil {
			return err
		}

		deleteResponse := apiPb.DeleteResponse{Result: true}

		response, err := json.Marshal(&deleteResponse)

		msg := &sarama.ProducerMessage{
			Topic:     "DeleteResponse",
			Partition: -1,
			Value:     sarama.ByteEncoder(response),
		}

		_, _, err = c.producer.SendMessage(msg)
	}

	return nil
}

func DeleteClaim(deleteHandler *DeleteHandler) {
	for {
		if err := deleteHandler.consumerGroup.Consume(context.Background(), []string{"DeleteRequest"}, deleteHandler); err != nil {
			log.Printf("on consume: %v", err)
			time.Sleep(time.Second * 10)
		}
	}
}
