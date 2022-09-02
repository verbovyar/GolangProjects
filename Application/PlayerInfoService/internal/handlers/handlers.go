package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-telegram-bot-api/telegram-bot-api/api/apiPb"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/domain"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/repositories/interfaces"
	"github.com/go-telegram-bot-api/telegram-bot-api/kafka"
	"log"
	"os"
	"os/signal"
)

func New(repository interfaces.Repository) *Handlers {
	return &Handlers{
		repository: repository,
		producer:   kafka.NewProducer(),
		consumer:   kafka.NewConsumer(),
	}
}

type Handlers struct {
	apiPb.UnimplementedPlayersServiceServer

	repository interfaces.Repository
	producer   sarama.SyncProducer
	consumer   sarama.Consumer
}

func (s *Handlers) List(ctx context.Context, in *apiPb.ListRequest) (*apiPb.ListResponse, error) {
	players := s.repository.List()
	playersDto := make([]*apiPb.ListResponse_Player, len(players))

	for i, player := range players {
		playersDto[i] = &apiPb.ListResponse_Player{
			Name:        player.GetName(),
			Club:        player.GetClub(),
			Id:          int32(player.GetId()),
			Nationality: player.GetNationality()}
	}

	response := apiPb.ListResponse{Players: playersDto}

	return &response, nil
}

func (s *Handlers) Add(ctx context.Context, in *apiPb.AddRequest) (*apiPb.AddResponse, error) {

	claim, err := s.consumer.ConsumePartition("AddRequest", 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		var addRequest *apiPb.AddRequest

		select {
		case err = <-claim.Errors():
			log.Println(err)
		case msg := <-claim.Messages():
			err = json.Unmarshal(msg.Value, addRequest)
			if err != nil {
				return nil, err
			}
		case <-signals:
			return nil, nil
		}

		player, err := domain.NewPlayer(addRequest.Name, addRequest.Club, addRequest.Nationality)
		if err != nil {
			return nil, err
		}

		err = s.repository.Add(player)
		if err != nil {
			return nil, err
		}
		addResponse := apiPb.AddResponse{Id: int32(player.Id)}

		response, err := json.Marshal(addResponse)

		msg := &sarama.ProducerMessage{
			Topic:     "AddResponse",
			Partition: -1,
			Value:     sarama.ByteEncoder(response),
		}

		_, _, err = s.producer.SendMessage(msg)
	}
}

func (s *Handlers) Update(ctx context.Context, in *apiPb.UpdateRequest) (*apiPb.UpdateResponse, error) {
	player, err := domain.NewPlayer(in.Name, in.Club, in.Nationality)
	if err != nil {
		return nil, err
	}

	err = s.repository.Update(player, uint(in.Id))
	if err != nil {
		return nil, err
	}
	player.Id = uint(in.Id)
	response := apiPb.UpdateResponse{Id: int32(uint(player.Id))}

	return &response, nil
}

func (s *Handlers) Delete(ctx context.Context, in *apiPb.DeleteRequest) (*apiPb.DeleteResponse, error) {
	err := s.repository.Delete(uint(in.Id))
	if err != nil {
		return nil, err
	}

	response := apiPb.DeleteResponse{Result: true}

	return &response, nil
}
