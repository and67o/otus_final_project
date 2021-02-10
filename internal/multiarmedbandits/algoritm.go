package multiarmedbandits

import (
	"math"

	"github.com/and67o/otus_project/internal/model"
)

func Get(banners []model.Banner) int64 {
	var selectedBannerID int64
	var count, maxCount float64
	bannersCount := len(banners)

	for _, banner := range banners {
		count = score(banner.CountClick, banner.CountShow, bannersCount)

		if count > maxCount || bannersCount == 0 {
			selectedBannerID = banner.IDBanner
			maxCount = count
		}
	}

	return selectedBannerID
}

func score(countClick int64, countShow int64, bannersCount int) float64 {
	value := float64(countClick) / float64(countShow)

	return value + math.Sqrt((2*math.Log(float64(bannersCount)))/float64(countShow))
}
