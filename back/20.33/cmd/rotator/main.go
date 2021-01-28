package main

import (
	"context"
	"flag"
	"github.com/and67o/otus_project/internal/app"
	"github.com/and67o/otus_project/internal/configuration"
	"github.com/and67o/otus_project/internal/logger"
	"github.com/and67o/otus_project/internal/server"
	"github.com/and67o/otus_project/internal/storage/sql"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := configuration.New(configFile)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := sql.New(config.DB)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(config.Logger)
	if err != nil {
		log.Fatal(err)
	}

	//queue, err := rmq.New(config.RabbitMQ, logg)
	//if err != nil {
	//	log.Fatal(err)
	//}

	rotator := app.New(storage, logg/*, queue*/)


	GRPCServer := server.New(rotator, config.Server)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}

	//go watchSignals(GRPCServer)

	logg.Info("calendar is running...")

	go func() {
		defer wg.Done()
		if err = GRPCServer.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
		}
	}()
	//go signalChan(logg, http, grpc)
}

//func signalChan(log *zap.Logger, srv ...srv.Stopper) {
//	signals := make(chan os.Signal, 1)
//	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
//	fmt.Printf("Got %v...\n", <-signals)
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
//	defer cancel()
//	for _, s := range srv {
//		err := s.Stop(ctx)
//		if err != nil {
//			log.Error("failed to stop", zap.Error(err))
//		}
//	}
//}

func watchSignals(server server.GRPC) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	<-signals
	signal.Stop(signals)

	err := server.Stop()
	if err != nil {
		log.Fatal(err)
	}
}

func startServer(wg *sync.WaitGroup, ctx context.Context, s server.GRPC, logg logger.Interface) {
	defer wg.Done()
	logg.Info("starting  server")

	err := s.Start(ctx)
	if err != nil {
		logg.Fatal("failed to start server: " + err.Error())
	}
}
