package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/and67o/otus_project/internal/app"
	"github.com/and67o/otus_project/internal/configuration"
	"github.com/and67o/otus_project/internal/logger"
	rmq "github.com/and67o/otus_project/internal/queue"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue, err := rmq.New(config.Rabbit, logg)
	if err != nil {
		log.Fatal(err)
	}

	rotator := app.New(storage, logg, queue)

	GRPCServer := server.New(rotator, config.Server)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-ctx.Done():
		case sig := <-signals:
			logg.Info(fmt.Sprintf("signal -  %s", sig))
			cancel()
		}
	}()
	wg := sync.WaitGroup{}

	//go watchSignals(GRPCServer)

	wg.Add(1)
	startServer(&wg, ctx, GRPCServer, logg)
	wg.Wait()
}

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
