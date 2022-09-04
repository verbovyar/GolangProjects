package updateHandler

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/go-telegram-bot-api/telegram-bot-api/api/apiPb"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/domain"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/kafkaHandlers/structs"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/repositories/interfaces"
	"log"
	"time"
)

type UpdateHandler struct {
	producer      sarama.SyncProducer
	repository    interfaces.Repository
	consumerGroup sarama.ConsumerGroup
}

func NewUpdateHandler(producer sarama.SyncProducer, repository interfaces.Repository, consumerGroup sarama.ConsumerGroup) *UpdateHandler {
	return &UpdateHandler{
		producer:      producer,
		repository:    repository,
		consumerGroup: consumerGroup,
	}
}

func (c *UpdateHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *UpdateHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *UpdateHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var updateRequest structs.UpdateRequest
		err := json.Unmarshal(msg.Value, &updateRequest)
		if err != nil {
			log.Printf("income data %v: %v", string(msg.Value), err)
			continue
		}

		player, err := domain.NewPlayer(updateRequest.Name, updateRequest.Club, updateRequest.Nationality)
		if err != nil {
			return err
		}

		err = c.repository.Update(player, uint(updateRequest.Id))
		if err != nil {
			return err
		}
		player.Id = uint(updateRequest.Id)
		listResponse := apiPb.UpdateResponse{Id: int32(player.Id)}

		response, err := json.Marshal(&listResponse)

		msg := &sarama.ProducerMessage{
			Topic:     "UpdateResponse",
			Partition: -1,
			Value:     sarama.ByteEncoder(response),
		}

		_, _, err = c.producer.SendMessage(msg)
	}

	return nil
}

func UpdateClaim(updateHandler *UpdateHandler) {
	for {
		if err := updateHandler.consumerGroup.Consume(context.Background(), []string{"UpdateRequest"}, updateHandler); err != nil {
			log.Printf("on consume: %v", err)
			time.Sleep(time.Second * 10)
		}
	}
}
