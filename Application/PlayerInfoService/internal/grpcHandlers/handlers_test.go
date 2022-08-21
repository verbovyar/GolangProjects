package grpcHandlers

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/api/apiPb"
	mock_grpcHandlers "github.com/go-telegram-bot-api/telegram-bot-api/internal/grpcHandlers/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPlayersInfoService_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlayer := mock_grpcHandlers.NewMockPlayersServiceClient(ctrl)

	in := &apiPb.AddRequest{
		Name:        "Verbov",
		Club:        "Barcelona",
		Nationality: "Russian",
	}

	var err error
	mockPlayer.EXPECT().Add(context.Background(), in).Return(nil, err)

	_, err = mockPlayer.Add(context.Background(), in)
	require.NoError(t, err)
}

func TestPlayersInfoService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlayer := mock_grpcHandlers.NewMockPlayersServiceClient(ctrl)

	in := &apiPb.DeleteRequest{
		Id: 24,
	}

	var err error
	mockPlayer.EXPECT().Delete(context.Background(), in).Return(nil, err)

	_, err = mockPlayer.Delete(context.Background(), in)
	require.NoError(t, err)
}

func TestPlayersInfoService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlayer := mock_grpcHandlers.NewMockPlayersServiceClient(ctrl)

	in := &apiPb.ListRequest{}

	var err error
	mockPlayer.EXPECT().List(context.Background(), in).Return(nil, err)

	_, err = mockPlayer.List(context.Background(), in)
	require.NoError(t, err)
}

func TestPlayersInfoService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlayer := mock_grpcHandlers.NewMockPlayersServiceClient(ctrl)

	in := &apiPb.UpdateRequest{
		Id:          24,
		Name:        "Verbov",
		Club:        "Psg",
		Nationality: "Russian",
	}

	var err error
	mockPlayer.EXPECT().Update(context.Background(), in).Return(nil, err)

	_, err = mockPlayer.Update(context.Background(), in)
	require.NoError(t, err)
}
