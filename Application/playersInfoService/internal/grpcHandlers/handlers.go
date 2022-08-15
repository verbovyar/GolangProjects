package grpcHandlers

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/domain"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/repositories/interfaces"
	pb "github.com/go-telegram-bot-api/telegram-bot-api/pkg/apiPb"
)

func New(repository interfaces.Repository) *Handlers {
	return &Handlers{
		repository: repository,
	}
}

type Handlers struct {
	pb.UnimplementedPlayersServiceServer

	repository interfaces.Repository
}

func (s *Handlers) List(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	players := s.repository.List()
	playersDto := make([]*pb.ListResponse_Player, len(players))

	for i, player := range players {
		playersDto[i] = &pb.ListResponse_Player{
			Name:        player.GetName(),
			Club:        player.GetClub(),
			Id:          int32(player.GetId()),
			Nationality: player.GetNationality()}
	}

	response := pb.ListResponse{Players: playersDto}

	return &response, nil
}

func (s *Handlers) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddResponse, error) {

	player, err := domain.NewPlayer(in.Name, in.Club, in.Nationality)
	if err != nil {
		return nil, err
	}

	err = s.repository.Add(player)
	if err != nil {
		return nil, err
	}

	response := pb.AddResponse{Id: int32(uint(player.Id))}

	return &response, nil
}

func (s *Handlers) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	player, err := domain.NewPlayer(in.Name, in.Club, in.Nationality)
	if err != nil {
		return nil, err
	}

	err = s.repository.Update(player, uint(in.Id))
	if err != nil {
		return nil, err
	}
	player.Id = uint(in.Id)
	response := pb.UpdateResponse{Id: int32(uint(player.Id))}

	return &response, nil
}

func (s *Handlers) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := s.repository.Delete(uint(in.Id))
	if err != nil {
		return nil, err
	}

	response := pb.DeleteResponse{Result: true}

	return &response, nil
}
