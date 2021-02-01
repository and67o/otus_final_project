package app

import (
	"context"
	"fmt"
	"time"

	"github.com/and67o/otus_project/internal/logger"
	"github.com/and67o/otus_project/internal/model"
	"github.com/and67o/otus_project/internal/multiarmedbandits"
	rmq "github.com/and67o/otus_project/internal/queue"
	server "github.com/and67o/otus_project/internal/server/pb"
	"github.com/and67o/otus_project/internal/storage/sql"
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
		BannerID: int(request.BannerId),
		SlotID:   int(request.SlotId),
	}

	err := a.Storage.AddBanner(&banner)
	if err != nil {
		return nil, fmt.Errorf("add banner: %w", err)
	}

	return &server.AddBannerResponse{}, nil
}

func (a *App) DeleteBanner(ctx context.Context, request *server.DeleteBannerRequest) (*server.DeleteBannerResponse, error) {
	banner := model.BannerPlace{
		BannerID: int(request.BannerId),
		SlotID:   int(request.SlotId),
	}

	err := a.Storage.DeleteBanner(&banner)
	if err != nil {
		return nil, fmt.Errorf("delete banner: %w", err)
	}

	return &server.DeleteBannerResponse{}, err
}

func (a *App) ClickBanner(ctx context.Context, request *server.ClickBannerRequest) (*server.ClickBannerResponse, error) {
	err := a.Storage.IncClickCount(request.SlotId, request.GroupId, request.BannerId)
	if err != nil {
		return nil, fmt.Errorf("click banner: %w", err)
	}

	err = a.Queue.Push(model.StatisticsEvent{
		Type:     model.TypeClick,
		IDSlot:   request.SlotId,
		IDBanner: request.BannerId,
		IDGroup:  request.GroupId,
		Date:     time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("push to queue: %w", err)
	}

	return &server.ClickBannerResponse{}, nil
}

func (a *App) ShowBanner(ctx context.Context, request *server.ShowBannerRequest) (*server.ShowBannerResponse, error) {
	banners, err := a.Storage.Banners(request.SlotId, request.GroupId)

	if err != nil {
		return nil, fmt.Errorf("show banner: %w", err)
	}

	showBannerID := multiarmedbandits.Get(banners)

	if showBannerID > 0 {
		err = a.Storage.IncShowCount(request.SlotId, request.GroupId, showBannerID)
		if err != nil {
			return nil, fmt.Errorf("increment count: %w", err)
		}

		err = a.Queue.Push(model.StatisticsEvent{
			Type:     model.TypeShow,
			IDSlot:   request.SlotId,
			IDBanner: showBannerID,
			IDGroup:  request.GroupId,
			Date:     time.Now(),
		})
		if err != nil {
			return nil, fmt.Errorf("push to queue: %w", err)
		}
	}

	return &server.ShowBannerResponse{BannerId: showBannerID}, err
}
