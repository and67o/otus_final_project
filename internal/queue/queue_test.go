package rmq

import (
	"fmt"
	"github.com/and67o/otus_project/internal/configuration"
	"github.com/and67o/otus_project/internal/logger"
	"github.com/and67o/otus_project/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var logConf = configuration.LoggerConf{
	Level:   "INFO",
	File:    "../../logs/log.log",
	IsProd:  false,
	TraceOn: false,
}

var rabbitConf = configuration.RabbitMQ{
	User: "guest",
	Pass: "guest",
	Host: "127.0.0.1",
	Port: 5672,
}

func TestRabbit(t *testing.T) {
	logg, err := logger.New(logConf)
	require.NoError(t, err)

	r, err := New(rabbitConf, logg)
	require.NoError(t, err)

	err = r.Push(model.StatisticsEvent{
		Type:     0,
		IDSlot:   0,
		IDBanner: 0,
		IDGroup:  0,
		Date:     time.Time{},
	})
	require.NoError(t, err)

	fmt.Println(err)
}
