package app

import (
	"context"
	"fmt"
	"github.com/and67o/otus_project/internal/logger"
	"github.com/and67o/otus_project/internal/model"
	"github.com/and67o/otus_project/internal/multiArmedBandits"
	rmq "github.com/and67o/otus_project/internal/queue"
	server "github.com/and67o/otus_project/internal/server/pb"
	"github.com/and67o/otus_project/internal/storage/sql"
	"time"
)

type App struct {
	Logger  logger.Interface
	Storage sql.StorageAction
	Queue   rmq.Queue

	server.UnimplementedBannerRotationServer
}

func New(storage sql.StorageAction, logger logger.Interface, queue rmq.Queue) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
		Queue:   queue,
	}
}

func (a *App) AddBanner(ctx context.Context, request *server.AddBannerRequest) (*server.AddBannerResponse, error) {
	banner := model.BannerPlace{
		BannerId: int(request.BannerId),
		SlotId:   int(request.SlotId),
	}

	err := a.Storage.AddBanner(&banner)
	if err != nil {
		return nil, err
	}

	return &server.AddBannerResponse{}, nil
}

func (a *App) DeleteBanner(ctx context.Context, request *server.DeleteBannerRequest) (*server.DeleteBannerResponse, error) {
	banner := model.BannerPlace{
		BannerId: int(request.BannerId),
		SlotId:   int(request.SlotId),
	}
	err := a.Storage.DeleteBanner(&banner)
	if err != nil {
		return nil, err
	}

	return &server.DeleteBannerResponse{}, err
}

func (a *App) ClickBanner(ctx context.Context, request *server.ClickBannerRequest) (*server.ClickBannerResponse, error) {
	var err error

	err = a.Storage.IncClickCount(request.SlotId, request.GroupId, request.BannerId)
	if err != nil {
		return nil, err
	}

	err = a.Queue.Push(model.StatisticsEvent{
		Type:     model.TypeClick,
		IDSlot:   request.SlotId,
		IDBanner: request.BannerId,
		IDGroup:  request.GroupId,
		Date:     time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &server.ClickBannerResponse{}, nil
}

func (a *App) ShowBanner(ctx context.Context, request *server.ShowBannerRequest) (*server.ShowBannerResponse, error) {
	var err error

	banners, err := a.Storage.Banners(request.SlotId, request.GroupId)

	if err != nil {
		return nil, err
	}

	showBannerId := multiArmedBandits.Get(banners)

	if showBannerId > 0 {
		err = a.Storage.IncShowCount(request.SlotId, request.GroupId, showBannerId)
		if err != nil {
			return nil, err
		}

		err = a.Queue.Push(model.StatisticsEvent{
			Type:     model.TypeShow,
			IDSlot:   request.SlotId,
			IDBanner: showBannerId,
			IDGroup:  request.GroupId,
			Date:     time.Now(),
		})
		fmt.Println(22,err)
		if err != nil {
			return nil, err
		}
	}

	return &server.ShowBannerResponse{BannerId: showBannerId}, err
}
