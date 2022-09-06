package app

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"modules/api/gateAwayApiPb"
	"modules/config"
	"modules/internal/handlers"
	pb "modules/internal/infrastructure/playersInfoServiceClient/api/pbGoFiles"
	"modules/kafka"
	"modules/pkg/logging"
	"net"
	"net/http"
	"os"
	"runtime/trace"
)

func Run(config config.Config) {
	//go counters.GetCounters()

	if err := trace.Start(os.Stderr); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	// create new channel of type int
	ch := make(chan int)

	// start new anonymous goroutine
	go func() {
		// send 42 to channel
		ch <- 42
	}()
	// read from channel
	<-ch

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

	ctx, task := trace.NewTask(context.Background(), "started run")
	defer task.End()

	go func() {
		trace.WithRegion(ctx, "Run Rest", func() {
			go runRest(config.HostGrpcPort, config.HostRestPort, logger)
		})
	}()

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
	ctx, cancel := context.WithCancel(context.Background())
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
