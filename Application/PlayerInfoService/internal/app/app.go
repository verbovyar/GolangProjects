package app

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/api/apiPb"
	"github.com/go-telegram-bot-api/telegram-bot-api/config"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/botService"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/handlers"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/repositories/interfaces"
	repo "github.com/go-telegram-bot-api/telegram-bot-api/internal/repositories/players/db"
	"github.com/go-telegram-bot-api/telegram-bot-api/pkg/postgres"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Run(config config.Config) {
	p := postgres.New(config.ConnectionString)
	playersRepo := repo.New(p.Pool)

	botService.New(playersRepo)
	go runBot(config.ApiKey)

	runGRPCServer(playersRepo, config.Network, config.Port)
}

func runBot(apiKey string) {
	log.Println("start cmd")
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
	handlers := handlers.New(repo)
	apiPb.RegisterPlayersServiceServer(grpcServer, handlers)

	listener, err := net.Listen(network, hostGrpcPort)
	if err != nil {
		log.Printf("%s", error(err))
	}

	if err = grpcServer.Serve(listener); err != nil {
		log.Panic(err)
	}
}
