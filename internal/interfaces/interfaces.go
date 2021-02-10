package interfaces

import (
	"context"

	"github.com/and67o/otus_project/internal/model"
)

type Queue interface {
	Push(event model.StatisticsEvent) error
	OpenChanel() error
	CloseChannel() error
	CloseConnection() error
}

type Storage interface {
	AddBanner(ctx context.Context, b *model.BannerPlace) error
	DeleteBanner(ctx context.Context, b *model.BannerPlace) error

	Banners(ctx context.Context, slotID int64, groupID int64) ([]model.Banner, error)
	IncShowCount(ctx context.Context, slotID int64, groupID int64, bannerID int64) error
	IncClickCount(ctx context.Context, slotID int64, groupID int64, bannerID int64) error

	AddStatistics(ctx context.Context, stat *model.Statistics) error
	DeleteStatistics(ctx context.Context, stat *model.Statistics) error
	GetStatistics(ctx context.Context, stat *model.Statistics) (*model.Statistics, error)

	UpdateStatus(ctx context.Context, status model.BannerStatus, b *model.BannerPlace) error
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Fatal(msg string)
	Warn(msg string)
}

type GRPC interface {
	Stop() error
	Start(ctx context.Context) error
}
