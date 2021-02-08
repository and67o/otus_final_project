package app

import (
	"context"
	"fmt"
	"github.com/and67o/otus_project/internal/interfaces"
	"time"

	"github.com/and67o/otus_project/internal/model"
	"github.com/and67o/otus_project/internal/multiarmedbandits"
	server "github.com/and67o/otus_project/internal/server/pb"
)

type App struct {
	logger  interfaces.Logger
	storage interfaces.Storage
	queue   interfaces.Queue

	server.UnimplementedBannerRotationServer
}

func New(storage interfaces.Storage, logger interfaces.Logger, queue interfaces.Queue) *App {
	return &App{
		logger:  logger,
		storage: storage,
		queue:   queue,
	}
}

func (a *App) AddBanner(_ context.Context, request *server.AddBannerRequest) (*server.AddBannerResponse, error) {
	banner := model.BannerPlace{
		BannerID: int(request.BannerId),
		SlotID:   int(request.SlotId),
	}

	err := a.storage.AddBanner(&banner)
	if err != nil {
		return nil, fmt.Errorf("add banner: %w", err)
	}

	return &server.AddBannerResponse{}, nil
}

func (a *App) DeleteBanner(_ context.Context, request *server.DeleteBannerRequest) (*server.DeleteBannerResponse, error) {
	banner := model.BannerPlace{
		BannerID: int(request.BannerId),
		SlotID:   int(request.SlotId),
	}

	err := a.storage.DeleteBanner(&banner)
	if err != nil {
		return nil, fmt.Errorf("delete banner: %w", err)
	}

	return &server.DeleteBannerResponse{}, err
}

func (a *App) ClickBanner(_ context.Context, request *server.ClickBannerRequest) (*server.ClickBannerResponse, error) {
	err := a.storage.IncClickCount(request.SlotId, request.GroupId, request.BannerId)
	if err != nil {
		return nil, fmt.Errorf("click banner: %w", err)
	}

	err = a.queue.Push(model.StatisticsEvent{
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

func (a *App) ShowBanner(_ context.Context, request *server.ShowBannerRequest) (*server.ShowBannerResponse, error) {
	banners, err := a.storage.Banners(request.SlotId, request.GroupId)

	if err != nil {
		return nil, fmt.Errorf("show banner: %w", err)
	}

	showBannerID := multiarmedbandits.Get(banners)

	if showBannerID > 0 {
		err = a.storage.IncShowCount(request.SlotId, request.GroupId, showBannerID)
		if err != nil {
			return nil, fmt.Errorf("increment count: %w", err)
		}

		err = a.queue.Push(model.StatisticsEvent{
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
