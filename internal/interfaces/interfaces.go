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
	AddBanner(b *model.BannerPlace) error
	DeleteBanner(b *model.BannerPlace) error

	Banners(slotID int64, groupID int64) ([]model.Banner, error)
	IncShowCount(slotID int64, groupID int64, bannerID int64) error
	IncClickCount(slotID int64, groupID int64, bannerID int64) error

	AddStatistics(stat *model.Statistics) error
	DeleteStatistics(stat *model.Statistics) error
	GetStatistics(stat *model.Statistics) (*model.Statistics, error)
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
