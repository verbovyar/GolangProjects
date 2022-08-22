package handlers

import (
	"context"
	"fmt"
	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
	"modules/api/gateAwayApiPb"
	pbGoFiles2 "modules/internal/infrastructure/playersInfoServiceClient/api/pbGoFiles"
	"modules/internal/utils"
)

func New(client pbGoFiles2.PlayersServiceClient) *Handlers {
	return &Handlers{
		client: client,
	}
}

type Handlers struct {
	gateAwayApiPb.UnsafePlayersInfoGateAwayServer
	client pbGoFiles2.PlayersServiceClient
}

func (h *Handlers) GetAll(ctx context.Context, in *gateAwayApiPb.GetAllRequest) (*gateAwayApiPb.GetAllResponse, error) {
	listRequest := &pbGoFiles2.ListRequest{}

	response, err := h.client.List(ctx, listRequest)
	if err != nil {
		fmt.Printf("list request error %v", err)
	}

	playersDto := make([]*gateAwayApiPb.GetAllResponse_Player, len(response.Players))

	for i, player := range response.Players {
		playersDto[i] = &gateAwayApiPb.GetAllResponse_Player{
			Name:        player.Name,
			Club:        player.Club,
			Id:          player.Id,
			Nationality: player.Nationality}
	}

	getAllResponse := gateAwayApiPb.GetAllResponse{Players: playersDto}

	return &getAllResponse, nil
}

func (h *Handlers) Post(ctx context.Context, in *gateAwayApiPb.PostRequest) (*gateAwayApiPb.PostResponse, error) {
	err := utils.ValidateAddRequest(in.Name, in.Club, in.Nationality)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	addRequest := &pbGoFiles2.AddRequest{
		Name:        in.Name,
		Club:        in.Club,
		Nationality: in.Nationality,
	}

	response, err := h.client.Add(ctx, addRequest)
	if err != nil {
		fmt.Printf("add erequest error %v", err)
	}

	postResponse := gateAwayApiPb.PostResponse{Id: response.Id}

	return &postResponse, nil
}

func (h *Handlers) Put(ctx context.Context, in *gateAwayApiPb.PutRequest) (*gateAwayApiPb.PutResponse, error) {

	err := utils.ValidateUpdateRequest(in.Name, in.Club, in.Nationality, in.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	updateRequest := &pbGoFiles2.UpdateRequest{
		Name:        in.Name,
		Club:        in.Club,
		Nationality: in.Nationality,
		Id:          in.Id,
	}

	response, err := h.client.Update(ctx, updateRequest)
	if err != nil {
		fmt.Printf("update request error %v", err)
	}

	putResponse := gateAwayApiPb.PutResponse{Id: response.Id}

	return &putResponse, nil
}

func (h *Handlers) Drop(ctx context.Context, in *gateAwayApiPb.DropRequest) (*gateAwayApiPb.DropResponse, error) {
	err := utils.ValidateDeleteRequest(in.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	deleteRequest := &pbGoFiles2.DeleteRequest{Id: in.Id}

	response, err := h.client.Delete(ctx, deleteRequest)
	if err != nil {
		fmt.Printf("delete request error %v", err)
	}

	deleteResponse := gateAwayApiPb.DropResponse{Result: response.Result}

	return &deleteResponse, nil
}
