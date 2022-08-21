package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"modules/api/gateAwayApiPb"
	"modules/config"
	pb "modules/infrastructure/playersInfoServiceClient/api/pbGoFiles"
	"modules/internal/grpcHandlers"
	"net"
	"net/http"
)

func main() {
	config, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("Started main")

	connect, err := grpc.Dial(config.PlayerInfoServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := pb.NewPlayersServiceClient(connect)

	go runRest(config.HostGrpcPort, config.HostRestPort)
	runGrpc(client, config.Network, config.HostGrpcPort)
}

func runGrpc(client pb.PlayersServiceClient, network, hostGrpcPort string) {
	fmt.Println("Started grpc")
	listener, err := net.Listen(network, hostGrpcPort)
	if err != nil {
		fmt.Printf("listen error %v\n", err)
	}

	handlers := grpcHandlers.New(client)

	grpcServer := grpc.NewServer()
	gateAwayApiPb.RegisterPlayersInfoGateAwayServer(grpcServer, handlers)

	err = grpcServer.Serve(listener)
	if err != nil {
		fmt.Printf("serve error%v\n", err)
	}
}

func runRest(hostGrpcPort, hostRestPort string) {
	fmt.Println("Started rest")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := gateAwayApiPb.RegisterPlayersInfoGateAwayHandlerFromEndpoint(ctx, mux, hostGrpcPort, opts); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(hostRestPort, mux); err != nil {
		panic(err)
	}
}
