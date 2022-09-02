package handlers

import (
	"context"
	"encoding/json"
	_ "expvar"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
	"log"
	"modules/api/gateAwayApiPb"
	pbGoFiles2 "modules/internal/infrastructure/playersInfoServiceClient/api/pbGoFiles"
	"modules/internal/utils"
	"modules/pkg/logging"
	"os"
	"os/signal"
)

func New(client pbGoFiles2.PlayersServiceClient, producer sarama.SyncProducer, logger logging.Logger) *Handlers {
	return &Handlers{
		client:   client,
		producer: producer,
		logger:   logger,
	}
}

type Handlers struct {
	gateAwayApiPb.UnimplementedPlayersInfoGateAwayServer
	client   pbGoFiles2.PlayersServiceClient
	producer sarama.SyncProducer
	consumer sarama.Consumer
	logger   logging.Logger
}

func (h *Handlers) GetAll(ctx context.Context, in *gateAwayApiPb.GetAllRequest) (*gateAwayApiPb.GetAllResponse, error) {
	listRequest := &pbGoFiles2.ListRequest{}
	h.logger.Info("Create List request from players info service")

	response, err := h.client.List(ctx, listRequest)
	h.logger.Info("Get List response")
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
	h.logger.Info("Overwriting data in GetAll response")

	getAllResponse := gateAwayApiPb.GetAllResponse{Players: playersDto}

	return &getAllResponse, nil
}

func (h *Handlers) Post(ctx context.Context, in *gateAwayApiPb.PostRequest) (*gateAwayApiPb.PostResponse, error) {
	err := utils.ValidateAddRequest(in.Name, in.Club, in.Nationality, h.logger)
	h.logger.Info("Validate Add request data")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request, err := json.Marshal(in)
	if err != nil {
		fmt.Print(err)
	}

	msg := &sarama.ProducerMessage{
		Topic:     "AddRequest",
		Partition: -1,
		Value:     sarama.ByteEncoder(request),
	}
	partition, offset, err := h.producer.SendMessage(msg)
	h.logger.Info("info about message( partition:%v, offset:%v, error:%v )", partition, offset, err)
	if err != nil {
		h.logger.Info("producer send msg error")
	}

	//TODO GetResponse from consumer
	claim, err := h.consumer.ConsumePartition("AddRequest", 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var response *pbGoFiles2.AddResponse
	select {
	case err = <-claim.Errors():
		log.Println(err)
	case msg := <-claim.Messages():
		err = json.Unmarshal(msg.Value, response)
		if err != nil {
			return nil, err
		}
	case <-signals:
		return nil, nil
	}

	postResponse := gateAwayApiPb.PostResponse{Id: response.Id}
	h.logger.Info("Get Post response")

	return &postResponse, nil
}

func (h *Handlers) Put(ctx context.Context, in *gateAwayApiPb.PutRequest) (*gateAwayApiPb.PutResponse, error) {
	err := utils.ValidateUpdateRequest(in.Name, in.Club, in.Nationality, in.Id, h.logger)
	h.logger.Info("Validate Update request data")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	updateRequest := &pbGoFiles2.UpdateRequest{
		Name:        in.Name,
		Club:        in.Club,
		Nationality: in.Nationality,
		Id:          in.Id,
	}
	h.logger.Info("Overwriting data in Update request")

	response, err := h.client.Update(ctx, updateRequest)
	h.logger.Info("Get Update response")
	if err != nil {
		fmt.Printf("update request error %v", err)
	}

	putResponse := gateAwayApiPb.PutResponse{Id: response.Id}

	return &putResponse, nil
}

func (h *Handlers) Drop(ctx context.Context, in *gateAwayApiPb.DropRequest) (*gateAwayApiPb.DropResponse, error) {
	err := utils.ValidateDeleteRequest(in.Id, h.logger)
	h.logger.Info("Validate Delete request data")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	deleteRequest := &pbGoFiles2.DeleteRequest{Id: in.Id}
	h.logger.Info("Overwriting data in Delete request")

	response, err := h.client.Delete(ctx, deleteRequest)
	h.logger.Info("Get Delete response")
	if err != nil {
		fmt.Printf("delete request error %v", err)
	}

	deleteResponse := gateAwayApiPb.DropResponse{Result: response.Result}

	return &deleteResponse, nil
}
