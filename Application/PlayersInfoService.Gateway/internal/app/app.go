package app

import (
	"context"
	"expvar"
	_ "expvar"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"modules/api/gateAwayApiPb"
	"modules/config"
	"modules/internal/handlers"
	pb "modules/internal/infrastructure/playersInfoServiceClient/api/pbGoFiles"
	"modules/pkg/logging"
	"net"
	"net/http"
	r "runtime"
	"time"
)

func Run(config config.Config) {
	counter()

	connect, err := grpc.Dial(config.PlayerInfoServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := pb.NewPlayersServiceClient(connect)

	logger := logging.NewLogger()
	logger.Info("Create new client")

	var brokers = []string{"127.0.0.1:9095"} //TODO move in config file
	logger.Info("Added brokers")
	conf := sarama.NewConfig()
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.RequiredAcks = sarama.WaitForLocal
	conf.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, conf)
	logger.Info("Create new producer")
	//---------------------

	h := handlers.New(client, producer, logger)
	logger.Info("Create new handlers")

	go runRest(config.HostGrpcPort, config.HostRestPort, logger)
	runGrpcServer(h, config.Network, config.HostGrpcPort, logger)
}

var (
	startTime = time.Now().UTC()
)

func goroutines() interface{} {
	return r.NumGoroutine()
}

func cpu() interface{} {
	return r.NumCPU()
}

func uptime() interface{} {
	return int64(time.Since(startTime))
}

func counter() {
	expvar.Publish("Goroutines", expvar.Func(goroutines))
	expvar.Publish("Uptime", expvar.Func(uptime))
	expvar.Publish("Cpu", expvar.Func(cpu))
	go func() {
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			fmt.Print("counter error", err)
		}
	}()
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
