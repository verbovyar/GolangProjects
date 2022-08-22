//go:build integration
// +build integration

package test

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/api/apiPb"
	utils "github.com/go-telegram-bot-api/telegram-bot-api/integration-test/config"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
)

func NewClient() apiPb.PlayersServiceClient {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("%v", err)
	}
	conn, err := grpc.Dial(config.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	Client := apiPb.NewPlayersServiceClient(conn)
	return Client
}

func TestAdd(t *testing.T) {
	in := &apiPb.AddRequest{
		Name:        "Test",
		Club:        "Psg",
		Nationality: "Russian",
	}

	client := NewClient()

	t.Run("AddPass", func(t *testing.T) {
		response, err := client.Add(context.Background(), in)

		if err == nil {
			listIn := &apiPb.ListRequest{}
			listResponse, err := client.List(context.Background(), listIn)

			if err == nil {
				for _, player := range listResponse.Players {
					if player.Name == in.Name && player.Id == response.Id {

						inDelete := &apiPb.DeleteRequest{Id: response.Id}
						_, _ = client.Delete(context.Background(), inDelete)

						assert.True(t, true)
						return
					}
				}
			}
		}
		assert.True(t, false)
	})
}

func TestUpdate(t *testing.T) {
	in := &apiPb.UpdateRequest{
		Name:        "Test",
		Club:        "Barcelona",
		Nationality: "Russian",
		Id:          0,
	}

	client := NewClient()

	inAdd := &apiPb.AddRequest{
		Name:        "Test",
		Club:        "Psg",
		Nationality: "Russian",
	}

	addResponse, _ := client.Add(context.Background(), inAdd)
	in.Id = addResponse.Id

	t.Run("UpdatePass", func(t *testing.T) {
		response, err := client.Update(context.Background(), in)

		if err == nil {
			listIn := &apiPb.ListRequest{}
			listResponse, err := client.List(context.Background(), listIn)

			if err == nil {
				for _, player := range listResponse.Players {
					if player.Club == in.Club && player.Id == response.Id {

						inDelete := &apiPb.DeleteRequest{Id: response.Id}
						_, _ = client.Delete(context.Background(), inDelete)

						assert.True(t, true)
						return
					}
				}
			}
		}
		assert.True(t, false)
	})
}

func TestList(t *testing.T) {
	in := &apiPb.ListRequest{}

	inAdd := &apiPb.AddRequest{
		Name:        "Test",
		Club:        "Psg",
		Nationality: "Russian",
	}

	client := NewClient()

	addResponse, err := client.Add(context.Background(), inAdd)

	if err == nil {
		t.Run("ListPass", func(t *testing.T) {
			response, err := client.List(context.Background(), in)
			if err == nil {
				if len(response.Players) == 6 {
					for _, player := range response.Players {
						if player.Id == addResponse.Id {
							inDelete := &apiPb.DeleteRequest{Id: addResponse.Id}
							_, _ = client.Delete(context.Background(), inDelete)

							assert.True(t, true)
							return
						}
					}
				}
			}
		})
	}
	assert.True(t, false)
}

func TestDelete(t *testing.T) {
	inAddFirstElement := &apiPb.AddRequest{
		Name:        "TestFirstElement",
		Club:        "Psg",
		Nationality: "Russian",
	}
	inAddSecondElement := &apiPb.AddRequest{
		Name:        "TestSecondElement",
		Club:        "Psg",
		Nationality: "Russian",
	}

	client := NewClient()

	addResponse, _ := client.Add(context.Background(), inAddFirstElement)
	_, _ = client.Add(context.Background(), inAddSecondElement)

	in := &apiPb.DeleteRequest{Id: addResponse.Id}

	t.Run("DeletePass", func(t *testing.T) {

		_, err := client.Delete(context.Background(), in)

		if err == nil {
			listIn := &apiPb.ListRequest{}
			listResponse, _ := client.List(context.Background(), listIn)

			for _, player := range listResponse.Players {
				if len(listResponse.Players) == 1 && player.Name == "TestSecondElement" {
					inDelete := &apiPb.DeleteRequest{Id: listResponse.Players[0].Id}
					_, _ = client.Delete(context.Background(), inDelete)

					assert.True(t, true)
					return
				}
			}
		}
		assert.True(t, true)
	})
}
