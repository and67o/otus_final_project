package model

type BannerStatus int

const (
	BannerStatusActive = iota
	BannerStatusDeleted
)

type Banner struct {
	IDBanner   int64
	IDSlot     int64
	CountClick int64
	CountShow  int64
	Status     BannerStatus
}
