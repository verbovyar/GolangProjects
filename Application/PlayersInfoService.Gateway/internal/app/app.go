package app

import (
	"context"
	"expvar"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"modules/api/gateAwayApiPb"
	"modules/config"
	"modules/internal/handlers"
	pb "modules/internal/infrastructure/playersInfoServiceClient/api/pbGoFiles"
	"modules/kafka"
	"modules/pkg/logging"
	"net"
	"net/http"
	r "runtime"
	"time"
)

func Run(config config.Config) {
	connect, err := grpc.Dial(config.PlayerInfoServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := pb.NewPlayersServiceClient(connect)

	logger := logging.NewLogger()
	logger.Info("Create new client")

	producer, err := kafka.NewProducer()
	logger.Info("Create new producer")

	h := handlers.New(client, producer, logger)
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

	in := &gateAwayApiPb.CountersRequest{}
	_, err = handlers.Counters(context.Background(), in)
	if err != nil {
		return
	}

	err = grpcServer.Serve(listener)
	logger.Info("Serve server")
	if err != nil {
		fmt.Printf("serve error%v\n", err)
	}
}

func counters(mux *http.ServeMux) {
	goroutines := func(w http.ResponseWriter, req *http.Request) {
		_, err := io.WriteString(w, string(rune(r.NumGoroutine())))
		if err != nil {
			return
		}
	}
	mux.HandleFunc("goroutines", goroutines)

	cpu := func(w http.ResponseWriter, req *http.Request) {
		_, err := io.WriteString(w, string(rune(r.NumCPU())))
		if err != nil {
			return
		}
	}
	mux.HandleFunc("cpu", cpu)

	startTime := time.Now().UTC()
	uptime := func(w http.ResponseWriter, req *http.Request) {
		_, err := io.WriteString(w, string(rune(int64(time.Since(startTime)))))
		if err != nil {
			return
		}
	}
	mux.HandleFunc("uptime", uptime)

	mux.Handle("/stats", expvar.Handler())

	http.ListenAndServe(":8080", mux)
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

	counters(http.NewServeMux())

	err = http.ListenAndServe(hostRestPort, mux)
	logger.Info("Listen and serve Rest")
	if err != nil {
		panic(err)
	}
}
