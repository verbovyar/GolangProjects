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
	pb "modules/internal/infrastructure/playersInfoServiceClient/api/pbGoFiles"
	"modules/internal/utils"
	"modules/pkg/logging"
	"time"
)

func New(client pb.PlayersServiceClient, producer sarama.SyncProducer, consumer sarama.Consumer, logger logging.Logger) *Handlers {
	return &Handlers{
		client:   client,
		producer: producer,
		logger:   logger,
		consumer: consumer,
	}
}

// TODO: Add logger in all functions
type Handlers struct {
	gateAwayApiPb.UnimplementedPlayersInfoGateAwayServer
	client   pb.PlayersServiceClient
	producer sarama.SyncProducer
	consumer sarama.Consumer
	logger   logging.Logger
}

type ListRequest struct {
}

func (h *Handlers) GetAll(ctx context.Context, in *gateAwayApiPb.GetAllRequest) (*gateAwayApiPb.GetAllResponse, error) {
	listRequest := &ListRequest{}

	request, err := json.Marshal(listRequest)
	if err != nil {
		fmt.Print(err)
	}

	msg := &sarama.ProducerMessage{
		Topic:     "ListRequest",
		Partition: -1,
		Value:     sarama.ByteEncoder(request),
	}
	partition, offset, err := h.producer.SendMessage(msg)
	h.logger.Info("info about message (partition:%v, offset:%v, error:%v)", partition, offset, err)
	if err != nil {
		h.logger.Info("producer send msg error")
	}

	time.Sleep(time.Second * 2)
	claim, err := h.consumer.ConsumePartition("ListResponse", 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Println(err)
	}

	var response pb.ListResponse
	for {
		select {
		case err = <-claim.Errors():
			log.Println(err)
		case msg := <-claim.Messages():
			err = json.Unmarshal(msg.Value, &response)
			if err != nil {
				return nil, err
			}
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
}

type AddRequest struct {
	Name        string `json:"name"`
	Club        string `json:"club"`
	Nationality string `json:"nationality"`
}

func (h *Handlers) Post(ctx context.Context, in *gateAwayApiPb.PostRequest) (*gateAwayApiPb.PostResponse, error) {
	err := utils.ValidateAddRequest(in.Name, in.Club, in.Nationality, h.logger)
	h.logger.Info("Validate Add request data")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	addRequest := &AddRequest{
		Name:        in.Name,
		Club:        in.Club,
		Nationality: in.Nationality,
	}

	request, err := json.Marshal(addRequest)
	if err != nil {
		fmt.Print(err)
	}

	msg := &sarama.ProducerMessage{
		Topic:     "AddRequest",
		Partition: -1,
		Value:     sarama.ByteEncoder(request),
	}
	partition, offset, err := h.producer.SendMessage(msg)
	h.logger.Info("info about message (partition:%v, offset:%v, error:%v)", partition, offset, err)
	if err != nil {
		h.logger.Info("producer send msg error")
	}

	time.Sleep(time.Second * 2)
	claim, err := h.consumer.ConsumePartition("AddResponse", 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Println(err)
	}

	var response pb.AddResponse
	for {
		select {
		case err = <-claim.Errors():
			log.Println(err)
		case msg := <-claim.Messages():
			err = json.Unmarshal(msg.Value, &response)
			if err != nil {
				return nil, err
			}
		}

		postResponse := &gateAwayApiPb.PostResponse{Id: response.Id}
		h.logger.Info("Get Post response")

		return postResponse, nil
	}
}

type UpdateRequest struct {
	Name        string `json:"name"`
	Club        string `json:"club"`
	Nationality string `json:"nationality"`
	Id          int32  `json:"id"`
}

func (h *Handlers) Put(ctx context.Context, in *gateAwayApiPb.PutRequest) (*gateAwayApiPb.PutResponse, error) {
	err := utils.ValidateUpdateRequest(in.Name, in.Club, in.Nationality, in.Id, h.logger)
	h.logger.Info("Validate Update request data")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	updateRequest := &UpdateRequest{
		Name:        in.Name,
		Club:        in.Club,
		Nationality: in.Nationality,
		Id:          in.Id,
	}
	h.logger.Info("Overwriting data in Update request")

	request, err := json.Marshal(updateRequest)

	msg := &sarama.ProducerMessage{
		Topic:     "UpdateRequest",
		Partition: -1,
		Value:     sarama.ByteEncoder(request),
	}
	partition, offset, err := h.producer.SendMessage(msg)
	h.logger.Info("info about message (partition:%v, offset:%v, error:%v)", partition, offset, err)
	if err != nil {
		h.logger.Info("producer send msg error")
	}

	time.Sleep(time.Second * 2)
	claim, err := h.consumer.ConsumePartition("UpdateResponse", 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Println(err)
	}

	var response pb.UpdateResponse
	for {
		select {
		case err = <-claim.Errors():
			log.Println(err)
		case msg := <-claim.Messages():
			err = json.Unmarshal(msg.Value, &response)
			if err != nil {
				return nil, err
			}
		}

		putResponse := gateAwayApiPb.PutResponse{Id: response.Id}

		return &putResponse, nil
	}
}

type DeleteRequest struct {
	Id int32 `json:"id"`
}

func (h *Handlers) Drop(ctx context.Context, in *gateAwayApiPb.DropRequest) (*gateAwayApiPb.DropResponse, error) {
	err := utils.ValidateDeleteRequest(in.Id, h.logger)
	h.logger.Info("Validate Delete request data")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	deleteRequest := &DeleteRequest{Id: in.Id}
	h.logger.Info("Overwriting data in Delete request")

	request, err := json.Marshal(deleteRequest)

	msg := &sarama.ProducerMessage{
		Topic:     "DeleteRequest",
		Partition: -1,
		Value:     sarama.ByteEncoder(request),
	}
	partition, offset, err := h.producer.SendMessage(msg)
	h.logger.Info("info about message (partition:%v, offset:%v, error:%v)", partition, offset, err)
	if err != nil {
		h.logger.Info("producer send msg error")
	}

	time.Sleep(time.Second * 2)
	claim, err := h.consumer.ConsumePartition("DeleteResponse", 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Println(err)
	}

	var response pb.DeleteResponse
	for {
		select {
		case err = <-claim.Errors():
			log.Println(err)
		case msg := <-claim.Messages():
			err = json.Unmarshal(msg.Value, &response)
			if err != nil {
				return nil, err
			}
		}

		deleteResponse := gateAwayApiPb.DropResponse{Result: response.Result}

		return &deleteResponse, nil
	}
}
