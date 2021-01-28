package model

const typeClick = 0
const typeShow = 1

type StatisticsEvent struct {
	Type     int
	IDSlot   int
	IDBanner int
	IDGroup  int
	Date     string
}
