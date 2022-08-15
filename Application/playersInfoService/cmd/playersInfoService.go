package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/botService"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/grpcHandlers"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/repositories/interfaces"
	repo "github.com/go-telegram-bot-api/telegram-bot-api/internal/repositories/players/db"
	"github.com/go-telegram-bot-api/telegram-bot-api/pkg/apiPb"
	"github.com/go-telegram-bot-api/telegram-bot-api/pkg/utils"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("%v", err)
	}

	playersRepo := repo.New(config.ConnectionString)

	botService.New(playersRepo)
	go runBot(config.ApiKey)

	runGRPCServer(playersRepo, config.Network, config.Port)
}

func runBot(apiKey string) {
	log.Println("start main")
	bot, err := botService.Init(apiKey)
	if err != nil {
		log.Panic(err)
	}

	botService.AddHandlers(bot)

	if err := bot.Run(); err != nil {
		log.Panic(err)
	}
}

func runGRPCServer(repo interfaces.Repository, network string, hostGrpcPort string) {
	grpcServer := grpc.NewServer()
	handlers := grpcHandlers.New(repo)
	apiPb.RegisterPlayersServiceServer(grpcServer, handlers)

	listener, err := net.Listen(network, hostGrpcPort)
	if err != nil {
		log.Printf("%s", error(err))
	}

	if err = grpcServer.Serve(listener); err != nil {
		log.Panic(err)
	}
}
