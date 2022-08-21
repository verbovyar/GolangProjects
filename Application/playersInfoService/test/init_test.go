//go:build integration
// +build integration

package test

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/api/apiPb"
	"github.com/go-telegram-bot-api/telegram-bot-api/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewClient() apiPb.PlayersServiceClient {
	config, err := config.LoadConfig(".")
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
