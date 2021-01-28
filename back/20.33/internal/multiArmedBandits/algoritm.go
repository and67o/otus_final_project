package multiArmedBandits

import (
	"github.com/and67o/otus_project/internal/model"
	"math"
)

func Get(banners []model.Banner) int64 {
	var selectedBannerId int64
	var count, maxCount float64
	bannersCount := len(banners)

	for _, banner := range banners {
		count = score(banner.ClickCount, banner.ShowCount, bannersCount)

		if count > maxCount || bannersCount == 0 {
			selectedBannerId = banner.ID
			maxCount = count
		}
	}

	return selectedBannerId
}

func score(countClick int64, countShow int64, bannersCount int) float64 {
	value := float64(countClick) / float64(countShow)
	return value + math.Sqrt((2*math.Log(float64(bannersCount)))/float64(countShow))
}
