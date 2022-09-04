package app

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"modules/api/gateAwayApiPb"
	"modules/config"
	"modules/internal/handlers"
	pb "modules/internal/infrastructure/playersInfoServiceClient/api/pbGoFiles"
	"modules/kafka"
	"modules/pkg/logging"
	"net"
	"net/http"
)

func Run(config config.Config) {
	//go counters.GetCounters()

	connect, err := grpc.Dial(config.PlayerInfoServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := pb.NewPlayersServiceClient(connect)

	logger := logging.NewLogger()
	logger.Info("Create new client")

	producer, err := kafka.NewProducer()
	logger.Info("Create new producer")

	consumer := kafka.NewConsumer()
	logger.Info("Create new consumer")

	h := handlers.New(client, producer, consumer, logger)
	logger.Info("Create new handlers")

	go runRest(config.HostGrpcPort, config.HostRestPort, logger)
	runGrpcServer(h, config.Network, config.HostGrpcPort, logger)
}

func runGrpcServer(handlers *handlers.Handlers, network, hostGrpcPort string, logger logging.Logger) {
	logger.Info("Started grpc")
	listener, err := net.Listen(network, hostGrpcPort)
	if err != nil {
		fmt.Printf("listen error %v\n", err)
	}
	logger.Info("Registered listener")

	grpcServer := grpc.NewServer()
	gateAwayApiPb.RegisterPlayersInfoGateAwayServer(grpcServer, handlers)
	logger.Info("Registered PlayersInfo Gate Away server")

	err = grpcServer.Serve(listener)
	logger.Info("Serve server")
	if err != nil {
		fmt.Printf("serve error%v\n", err)
	}
}

func runRest(hostGrpcPort, hostRestPort string, logger logging.Logger) {
	logger.Info("Started Rest")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gateAwayApiPb.RegisterPlayersInfoGateAwayHandlerFromEndpoint(ctx, mux, hostGrpcPort, opts)
	logger.Info("Registered players info gate away from endpoint")
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(hostRestPort, mux)
	logger.Info("Listen and serve Rest")
	if err != nil {
		panic(err)
	}
}
