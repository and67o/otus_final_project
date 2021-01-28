package app

import (
	"github.com/and67o/otus_project/internal/logger"
	rmq "github.com/and67o/otus_project/internal/queue"
	"github.com/and67o/otus_project/internal/storage/sql"
)


type App struct {
	Logger  logger.Interface
	Storage sql.StorageAction
	Queue rmq.Queue
}

func New(storage sql.StorageAction, logger logger.Interface, queue rmq.Queue) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
		Queue: queue,
	}
}
