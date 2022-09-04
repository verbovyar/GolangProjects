package listHandler

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

type ListHandler struct {
	producer      sarama.SyncProducer
	repository    interfaces.Repository
	consumerGroup sarama.ConsumerGroup
}

func NewListHandler(producer sarama.SyncProducer, repository interfaces.Repository, consumerGroup sarama.ConsumerGroup) *ListHandler {
	return &ListHandler{
		producer:      producer,
		repository:    repository,
		consumerGroup: consumerGroup,
	}
}

func (c *ListHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ListHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ListHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var listRequest structs.ListRequest
		err := json.Unmarshal(msg.Value, &listRequest)
		if err != nil {
			log.Printf("income data %v: %v", string(msg.Value), err)
			continue
		}

		players := c.repository.List()
		playersDto := make([]*apiPb.ListResponse_Player, len(players))

		for i, player := range players {
			playersDto[i] = &apiPb.ListResponse_Player{
				Name:        player.GetName(),
				Club:        player.GetClub(),
				Id:          int32(player.GetId()),
				Nationality: player.GetNationality()}
		}

		listResponse := apiPb.ListResponse{Players: playersDto}

		response, err := json.Marshal(&listResponse)

		msg := &sarama.ProducerMessage{
			Topic:     "ListResponse",
			Partition: -1,
			Value:     sarama.ByteEncoder(response),
		}

		_, _, err = c.producer.SendMessage(msg)
	}

	return nil
}

func ListClaim(listHandler *ListHandler) {
	for {
		if err := listHandler.consumerGroup.Consume(context.Background(), []string{"ListRequest"}, listHandler); err != nil {
			log.Printf("on consume: %v", err)
			time.Sleep(time.Second * 10)
		}
	}
}
