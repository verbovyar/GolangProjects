package handlers

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"modules/api/gateAwayApiPb"
	mock "modules/internal/handlers/mock"
	"modules/internal/infrastructure/playersInfoServiceClient/api/pbGoFiles"
	"testing"
)

func TestHandlers_Post(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockPlayersServiceClient(ctrl)

	addResponse := &pbGoFiles.AddResponse{Id: 2}

	client.EXPECT().Add(context.Background(), nil).Return(addResponse, nil)

	postRequest := &gateAwayApiPb.PostRequest{
		Name:        "Jane",
		Club:        "Liverpool",
		Nationality: "Argentinian",
	}

	handler := Handlers{}
	_, err := handler.Post(context.Background(), postRequest)

	require.NoError(t, err)
}

func TestHandlers_Put(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockPlayersServiceClient(ctrl)

	updateResponse := &pbGoFiles.UpdateResponse{Id: 2}

	client.EXPECT().Update(context.Background(), nil).Return(updateResponse, nil)

	putRequest := &gateAwayApiPb.PutRequest{
		Name:        "Jane",
		Club:        "Liverpool",
		Nationality: "Russian",
		Id:          2,
	}

	handler := Handlers{}
	_, err := handler.Put(context.Background(), putRequest)

	require.NoError(t, err)
}

func TestHandlers_Drop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockPlayersServiceClient(ctrl)

	deleteResponse := &pbGoFiles.DeleteResponse{Result: true}

	client.EXPECT().Delete(context.Background(), nil).Return(deleteResponse, nil)

	dropRequest := &gateAwayApiPb.DropRequest{Id: 2}

	handler := Handlers{}
	_, err := handler.Drop(context.Background(), dropRequest)

	require.NoError(t, err)
}
