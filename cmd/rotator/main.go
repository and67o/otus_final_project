package main

import (
	"context"
	"flag"
	"github.com/and67o/otus_project/internal/app"
	"github.com/and67o/otus_project/internal/configuration"
	"github.com/and67o/otus_project/internal/logger"
	rmq "github.com/and67o/otus_project/internal/queue"
	"github.com/and67o/otus_project/internal/server"
	"github.com/and67o/otus_project/internal/storage/sql"
	"log"
	"os"
	"os/signal"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue, err := rmq.New(config.Rabbit, logg)
	if err != nil {
		log.Fatal(err)
	}

	rotator := app.New(storage, logg, queue)

	GRPCServer := server.New(rotator, config.Server)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		<-signals
		signal.Stop(signals)
		cancel()

		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		//defer cancel()


		if err := GRPCServer.Stop(); err != nil {
			logg.Error("stop server"  + err.Error())
		}
	}()

	logg.Info("starting  server")
	err = GRPCServer.Start(ctx)
	if err != nil {
		logg.Fatal("failed to start server: " + err.Error())
	}
}
