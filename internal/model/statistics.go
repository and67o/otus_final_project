package model

import "time"

type EventType int

const (
	TypeClick EventType = iota
	TypeShow
)

type StatisticsEvent struct {
	Type     EventType
	IDSlot   int64
	IDBanner int64
	IDGroup  int64
	Date     time.Time
}

type Statistics struct {
	IDSlot     int64
	IDBanner   int64
	IDGroup    int64
	CountClick int64
	CountShow  int64
}
