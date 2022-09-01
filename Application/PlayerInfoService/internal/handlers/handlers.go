package handlers

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/go-telegram-bot-api/telegram-bot-api/api/apiPb"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/domain"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/repositories/interfaces"
	"log"
)

func New(repository interfaces.Repository) *Handlers {
	return &Handlers{
		repository: repository,
	}
}

type Handlers struct {
	apiPb.UnimplementedPlayersServiceServer

	repository interfaces.Repository
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
	var claim sarama.ConsumerGroupClaim

	for msg := range claim.Messages() {
		var addRequest *apiPb.AddRequest
		err := json.Unmarshal(msg.Value, addRequest)
		if err != nil {
			log.Printf("Unmarshal error: ^%v", err.Error())
			continue
		}

		player, err := domain.NewPlayer(in.Name, in.Club, in.Nationality)
		if err != nil {
			return nil, err
		}

		err = s.repository.Add(player)
		if err != nil {
			return nil, err
		}

		response := apiPb.AddResponse{Id: int32(uint(player.Id))}

		addResponseMarshal, err := json.Marshal(response)
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic:     "addRequest",
			Partition: -1,
			Value:     sarama.ByteEncoder(addResponseMarshal),
		})
		if err != nil {
			log.Printf("Cant pay order: %v", err)
		}

	}

	return &response, nil
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
