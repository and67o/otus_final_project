package model

import "time"

type EventType int

const TypeClick = 0
const TypeShow = 1

type StatisticsEvent struct {
	Type     int
	IDSlot   int64
	IDBanner int64
	IDGroup  int64
	Date     time.Time
}
