package server

import (
	"context"
	"fmt"
	"github.com/and67o/otus_project/internal/app"
	"github.com/and67o/otus_project/internal/configuration"
	storage2 "github.com/and67o/otus_project/internal/storage"
	"testing"
)

var configGrpc = configuration.GRPCConf{
	Host: "127.0.0.1",
	Port: "50051",
}

var storage = configuration.DBConf{
	User:   "admin",
	Pass:   "123",
	DBName: "go_api",
	Host:   "127.0.0.1",
	Port:   3306,
}

func TestRotatorServer_AddBannerToSlot(t *testing.T) {
	//t.Skip(true)
	storage, err := storage2.New(storage)

	a := app.New(storage, nil, nil)

	s := New(a, configGrpc)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = s.Start(ctx)
	fmt.Println(err)
}
