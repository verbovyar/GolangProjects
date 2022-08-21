package grpcHandlers

import (
	"context"
	"fmt"
	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
	pb "modules/infrastructure/playersInfoServiceClient/pbGoFiles"
	api "modules/pkg/gateAwayApiPb"
)

func New(client pb.PlayersServiceClient) *Handlers {
	return &Handlers{
		client: client,
	}
}

type Handlers struct {
	api.UnsafePlayersInfoGateAwayServer
	client pb.PlayersServiceClient
}

func (h *Handlers) GetAll(ctx context.Context, in *api.GetAllRequest) (*api.GetAllResponse, error) {
	listRequest := pb.ListRequest{}

	response, err := h.client.List(ctx, &listRequest)
	if err != nil {
		fmt.Printf("list request error %v", err)
	}

	playersDto := make([]*api.GetAllResponse_Player, len(response.Players))

	for i, player := range response.Players {
		playersDto[i] = &api.GetAllResponse_Player{
			Name:        player.Name,
			Club:        player.Club,
			Id:          player.Id,
			Nationality: player.Nationality}
	}

	getAllResponse := api.GetAllResponse{Players: playersDto}

	return &getAllResponse, nil
}

func (h *Handlers) Post(ctx context.Context, in *api.PostRequest) (*api.PostResponse, error) {
	err := validateAddRequest(in.Name, in.Club, in.Nationality)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	addRequest := pb.AddRequest{
		Name:        in.Name,
		Club:        in.Club,
		Nationality: in.Nationality,
	}

	response, err := h.client.Add(ctx, &addRequest)
	if err != nil {
		fmt.Printf("add erequest error %v", err)
	}

	postResponse := api.PostResponse{Id: response.Id}

	return &postResponse, nil
}

func (h *Handlers) Put(ctx context.Context, in *api.PutRequest) (*api.PutResponse, error) {

	err := validateUpdateRequest(in.Name, in.Club, in.Nationality, in.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	updateRequest := pb.UpdateRequest{
		Name:        in.Name,
		Club:        in.Club,
		Nationality: in.Nationality,
		Id:          in.Id,
	}

	response, err := h.client.Update(ctx, &updateRequest)
	if err != nil {
		fmt.Printf("update request error %v", err)
	}

	putResponse := api.PutResponse{Id: response.Id}

	return &putResponse, nil
}

func (h *Handlers) Drop(ctx context.Context, in *api.DropRequest) (*api.DropResponse, error) {
	err := validateDeleteRequest(in.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	deleteRequest := pb.DeleteRequest{Id: in.Id}

	response, err := h.client.Delete(ctx, &deleteRequest)
	if err != nil {
		fmt.Printf("delete request error %v", err)
	}

	deleteResponse := api.DropResponse{Result: response.Result}

	return &deleteResponse, nil
}
