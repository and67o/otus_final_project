package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/and67o/otus_project/internal/app"
	"github.com/and67o/otus_project/internal/configuration"
	"github.com/and67o/otus_project/internal/interfaces"
	"github.com/and67o/otus_project/internal/logger"
	rmq "github.com/and67o/otus_project/internal/queue"
	"github.com/and67o/otus_project/internal/server"
	storage2 "github.com/and67o/otus_project/internal/storage"
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

	storage, err := storage2.New(config.DB)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(config.Logger)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue, err := rmq.New(config.Rabbit)
	if err != nil {
		log.Fatal(err) // nolint: gocritic
	}

	err = queue.OpenChanel()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = queue.CloseChannel()
		log.Fatal(err)
	}()

	rotator := app.New(storage, logg, queue)

	GRPCServer := server.New(rotator, config.Server)

	go watchSignals(cancel, GRPCServer, logg, queue)

	logg.Info("starting  server")
	err = GRPCServer.Start(ctx)
	if err != nil {
		logg.Fatal("failed to start server: " + err.Error())
	}
}

func watchSignals(cancel context.CancelFunc, grpc interfaces.GRPC, logg interfaces.Logger, queue interfaces.Queue) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	<-signals
	signal.Stop(signals)
	cancel()

	err := grpc.Stop()
	if err != nil {
		logg.Error("stop server: " + err.Error())
	}

	err = queue.CloseConnection()
	if err != nil {
		logg.Error("queue error: " + err.Error())
	}
}
