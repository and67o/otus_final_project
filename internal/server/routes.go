package server

import (
	"github.com/and67o/otus_project/internal/model"
	"github.com/and67o/otus_project/internal/multiArmedBandits"
)

type GRPCRoutes interface {
	AddBanner(request *AddBannerRequest) error
	DeleteBanner(request *DeleteBannerRequest) error
	ClickBanner(request *ClickBannerRequest) error
	ShowBanner(request *ShowBannerRequest) (int, error)
}

func (s *Server) AddBanner(request *AddBannerRequest) error {
	banner := model.BannerPlace{
		BannerId: int(request.BannerId),
		SlotId:   int(request.SlotId),
	}

	err := s.app.Storage.AddBanner(&banner)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) DeleteBanner(request *DeleteBannerRequest) error {
	banner := model.BannerPlace{
		BannerId: int(request.BannerId),
		SlotId:   int(request.SlotId),
	}
	err := s.app.Storage.DeleteBanner(&banner)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) ClickBanner(request *ClickBannerRequest) error {
	panic("implement me")
}

func (s *Server) ShowBanner(request *ShowBannerRequest) (int, error) {
	banners, err := s.app.Storage.Banners(int(request.SlotId), int(request.GroupId))
	if err != nil {
		return 0, err
	}

	showBannerId := multiArmedBandits.Get(banners)
	if showBannerId > 0 {
		err = s.app.Storage.IncShowCount(request.SlotId, request.GroupId, showBannerId)
	}
	return 0, err
}
