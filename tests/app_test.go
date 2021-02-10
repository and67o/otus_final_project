package tests

import (
	"context"
	"github.com/and67o/otus_project/internal/app"
	"github.com/and67o/otus_project/internal/configuration"
	"github.com/and67o/otus_project/internal/interfaces"
	"github.com/and67o/otus_project/internal/logger"
	"github.com/and67o/otus_project/internal/model"
	rmq "github.com/and67o/otus_project/internal/queue"
	pb "github.com/and67o/otus_project/internal/server/pb"
	storage2 "github.com/and67o/otus_project/internal/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestCase struct {
	app     *app.App
	storage interfaces.Storage
	queue   interfaces.Queue
}

var exampleStatistics = model.Statistics{
	IDSlot:   1,
	IDBanner: 3,
	IDGroup:  1,
}

var configStorage = configuration.DBConf{
	User:   "admin",
	Pass:   "123",
	DBName: "go_api",
	Host:   "127.0.0.1",
	Port:   3306,
}

var configLogger = configuration.LoggerConf{
	Level:   "DEBUG",
	File:    "./testdata/log.log",
	IsProd:  false,
	TraceOn: false,
}

var rabbitConf = configuration.RabbitMQ{
	User:       "guest",
	Pass:       "guest",
	Host:       "127.0.0.1",
	Port:       5672,
	Durable:    true,
	AutoDelete: false,
	NoWait:     false,
	Internal:   false,
	Name:       "banner_exchange_queue",
	Kind:       "direct",
	Key:        "banner",
}

func TestApp(t *testing.T) {
	ts := getApp(t)

	ts.appUp(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func() {
		_ = ts.queue.CloseChannel()
		_ = ts.queue.CloseChannel()
	}()

	ts.appTest(ctx, t)

	ts.appDown(t)
}

func getApp(t *testing.T) *TestCase {
	ts := TestCase{}

	storage, err := storage2.New(configStorage)
	require.NoError(t, err)

	logg, err := logger.New(configLogger)
	require.NoError(t, err)

	queue, err := rmq.New(rabbitConf)
	require.NoError(t, err)
	err = queue.OpenChanel()
	require.NoError(t, err)

	ts.app = app.New(storage, logg, queue)
	ts.storage = storage
	ts.queue = queue

	return &ts
}

func (ts *TestCase) appUp(t *testing.T) {
	BannersUp(t, ts.storage)
	StatisticsUp(t, ts.storage)
}

func getTestBanners() []model.BannerPlace {
	return []model.BannerPlace{
		{BannerID: 1, SlotID: 2},
		{BannerID: 1, SlotID: 3},
		{BannerID: 2, SlotID: 1},
		{BannerID: 2, SlotID: 2},
		{BannerID: 2, SlotID: 3},
		{BannerID: 3, SlotID: 1},
		{BannerID: 4, SlotID: 3},
	}
}

func BannersUp(t *testing.T, storage interfaces.Storage) {
	banners := getTestBanners()
	for i := range getTestBanners() {
		err := storage.AddBanner(&banners[i])
		require.NoError(t, err)
	}
}

func StatisticsUp(t *testing.T, storage interfaces.Storage) {
	statistics := getTestStatistics()
	for i := range statistics {
		err := storage.AddStatistics(&statistics[i])
		require.NoError(t, err)
	}
}

func (ts *TestCase) appTest(ctx context.Context, t *testing.T) {
	var err error

	// уже есть такой баннер
	_, err = ts.app.AddBanner(ctx, &pb.AddBannerRequest{SlotId: 1, BannerId: 2})
	require.Error(t, err)

	// создать новый баннер
	_, err = ts.app.AddBanner(ctx, &pb.AddBannerRequest{SlotId: 1, BannerId: 1})
	require.NoError(t, err)

	// удалить баннер
	_, err = ts.app.DeleteBanner(ctx, &pb.DeleteBannerRequest{SlotId: 1, BannerId: 1})
	require.NoError(t, err)

	statBefore, err := ts.storage.GetStatistics(&exampleStatistics)
	require.NoError(t, err)

	// нажать на баннер
	_, err = ts.app.ClickBanner(ctx, &pb.ClickBannerRequest{
		SlotId:   1,
		BannerId: 3,
		GroupId:  1,
	})
	require.NoError(t, err)

	// проверить клик
	statAfterClick, err := ts.storage.GetStatistics(&exampleStatistics)
	require.NoError(t, err)

	require.Equal(t, statBefore.IDBanner, statAfterClick.IDBanner)
	require.Equal(t, statBefore.IDGroup, statAfterClick.IDGroup)
	require.Equal(t, statBefore.IDSlot, statAfterClick.IDSlot)
	require.Equal(t, statBefore.CountClick+1, statAfterClick.CountClick)

	response, err := ts.app.ShowBanner(ctx, &pb.ShowBannerRequest{
		SlotId:  1,
		GroupId: 1,
	})

	require.NoError(t, err)
	require.Equal(t, response.BannerId, int64(3))

	// проверить показы
	statAfterShow, err := ts.storage.GetStatistics(&exampleStatistics)
	require.NoError(t, err)

	require.Equal(t, statBefore.CountShow+1, statAfterShow.CountShow)
}

func (ts *TestCase) appDown(t *testing.T) {
	BannerDown(t, ts.storage)
	StatisticsDown(t, ts.storage)
}

func BannerDown(t *testing.T, storage interfaces.Storage) {
	banners := getTestBanners()
	for i := range getTestBanners() {
		err := storage.DeleteBanner(&banners[i])
		require.NoError(t, err)
	}
}

func getTestStatistics() []model.Statistics {
	return []model.Statistics{
		{IDSlot: 1, IDBanner: 3, IDGroup: 1, CountClick: 2, CountShow: 9},
		{IDSlot: 1, IDBanner: 3, IDGroup: 2, CountClick: 0, CountShow: 0},
		{IDSlot: 3, IDBanner: 2, IDGroup: 1, CountClick: 85, CountShow: 119},
		{IDSlot: 3, IDBanner: 4, IDGroup: 1, CountClick: 41, CountShow: 95},
		{IDSlot: 3, IDBanner: 1, IDGroup: 2, CountClick: 78, CountShow: 85},
		{IDSlot: 3, IDBanner: 2, IDGroup: 2, CountClick: 0, CountShow: 0},
		{IDSlot: 3, IDBanner: 4, IDGroup: 2, CountClick: 78, CountShow: 112},
		{IDSlot: 2, IDBanner: 1, IDGroup: 1, CountClick: 5, CountShow: 14},
		{IDSlot: 2, IDBanner: 2, IDGroup: 1, CountClick: 76, CountShow: 129},
		{IDSlot: 2, IDBanner: 1, IDGroup: 2, CountClick: 81, CountShow: 99},
		{IDSlot: 2, IDBanner: 2, IDGroup: 2, CountClick: 45, CountShow: 98},
		{IDSlot: 3, IDBanner: 1, IDGroup: 1, CountClick: 0, CountShow: 5},
	}
}

func StatisticsDown(t *testing.T, storage interfaces.Storage) {
	statistics := getTestStatistics()
	for i := range statistics {
		err := storage.DeleteStatistics(&statistics[i])
		require.NoError(t, err)
	}
}
