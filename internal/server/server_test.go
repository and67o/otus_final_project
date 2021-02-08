package server

import (
	"context"
	"fmt"
	"github.com/and67o/otus_project/internal/app"
	"github.com/and67o/otus_project/internal/configuration"
	"testing"
)

func TestRotatorServer_AddBannerToSlot(t *testing.T) {
	t.Skip(true)
	a := app.App{}
	config := configuration.GRPCConf{
		Host: "127.0.0.1",
		Port: "50051",
	}
	s := New(&a, config)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err:= s.Start(ctx)
	fmt.Println(err)
}
