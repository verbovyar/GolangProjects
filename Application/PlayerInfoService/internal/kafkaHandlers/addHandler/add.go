package addHandler

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

type AddHandler struct {
	producer      sarama.SyncProducer
	repository    interfaces.Repository
	consumerGroup sarama.ConsumerGroup
}

func NewAddHandler(producer sarama.SyncProducer, repository interfaces.Repository, consumerGroup sarama.ConsumerGroup) *AddHandler {
	return &AddHandler{
		producer:      producer,
		repository:    repository,
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
		var addRequest *structs.AddRequest
		err := json.Unmarshal(msg.Value, addRequest)
		if err != nil {
			log.Printf("income data %v: %v", string(msg.Value), err)
			continue
		}

		player, err := domain.NewPlayer(addRequest.Name, addRequest.Club, addRequest.Nationality)
		if err != nil {
			return err
		}

		err = c.repository.Add(player)
		if err != nil {
			return err
		}

		addResponse := apiPb.AddResponse{Id: int32(player.Id)}

		response, err := json.Marshal(addResponse)

		msg := &sarama.ProducerMessage{
			Topic:     "AddResponse",
			Partition: -1,
			Value:     sarama.ByteEncoder(response),
		}

		_, _, err = c.producer.SendMessage(msg)
	}

	return nil
}

func AddClaim(addHandler *AddHandler) {
	for {
		if err := addHandler.consumerGroup.Consume(context.Background(), []string{"AddRequest"}, addHandler); err != nil {
			log.Printf("on consume: %v", err)
			time.Sleep(time.Second * 10)
		}
	}
}
