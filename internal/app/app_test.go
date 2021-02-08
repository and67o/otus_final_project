package app

import (
	"context"
	"github.com/and67o/otus_project/internal/configuration"
	"github.com/and67o/otus_project/internal/logger"
	rmq "github.com/and67o/otus_project/internal/queue"
	server "github.com/and67o/otus_project/internal/server/pb"
	"github.com/and67o/otus_project/internal/storage/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

var dbConfig = configuration.DBConf{
	User:   "admin",
	Pass:   "123",
	DBName: "go_api",
	Host:   "localhost",
	Port:   3306,
}

var logConf = configuration.LoggerConf{
	Level:   "INFO",
	File:    "../../logs/log.log",
	IsProd:  false,
	TraceOn: false,
}

var rabbitConf = configuration.RabbitMQ{
	User: "guest",
	Pass: "guest",
	Host: "127.0.0.1",
	Port: 5672,
}

func TestBaseAction(t *testing.T) {
	db, err := sql.New(dbConfig)
	require.NoError(t, err)

	logg, err := logger.New(logConf)
	require.NoError(t, err)

	r, err := rmq.New(rabbitConf, logg)
	require.NoError(t, err)

	app := App{
		logger:                            logg,
		storage:                           db,
		queue:                             r,
		UnimplementedBannerRotationServer: server.UnimplementedBannerRotationServer{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err = app.DeleteBanner(ctx, &server.DeleteBannerRequest{
		BannerId: 2,
		SlotId:   1,
	})

	_, err = app.AddBanner(ctx, &server.AddBannerRequest{
		SlotId:   1,
		BannerId: 2,
	})
	require.NoError(t, err)


	require.NoError(t, err)

	_, err = app.ClickBanner(ctx, &server.ClickBannerRequest{
		SlotId:   2,
		BannerId: 1,
		GroupId:  2,
	})
	require.NoError(t, err)

	bannerResponse, err := app.ShowBanner(ctx, &server.ShowBannerRequest{
		SlotId:  2,
		GroupId: 1,
	})
	require.NoError(t, err)
	require.Greater(t, bannerResponse.BannerId, 0)
}
