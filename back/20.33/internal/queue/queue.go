package rmq

import (
	"context"
	"encoding/json"
	"github.com/and67o/otus_project/internal/configuration"
	"github.com/and67o/otus_project/internal/logger"
	"github.com/and67o/otus_project/internal/model"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	logger     logger.Interface
	connection *amqp.Connection
	channel    *amqp.Channel
}

type Queue interface {
}

const nameQueue = "banner_queue"

func New(config configuration.RabbitMQ, logg logger.Interface) (*RabbitMQ, error) {
	var err error

	res := &RabbitMQ{
		logger: logg,
	}

	connectionCtx, connectionCtxCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer connectionCtxCancel()

	res.connection, err = connectionWait(connectionCtx, getUrl(config))
	if err != nil {
		res.logger.Fatal("fail RabbitMQ connection")
		return nil, err
	}

	res.channel, err = res.connection.Channel()
	if err != nil {
		res.logger.Fatal("fail to open channel for RabbitMQ")
		return nil, err
	}
	_, err = res.channel.QueueDeclare(
		nameQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		res.logger.Fatal("fail create queue for RabbitMQ")
		return nil, err
	}

	return res, nil
}

func connectionWait(ctx context.Context, dsn string) (*amqp.Connection, error) {
	var err error

	var res *amqp.Connection

	for {
		res, err = amqp.Dial(dsn)
		if err == nil || ctx.Err() != nil {
			break
		}

		time.Sleep(time.Second)
	}

	return res, err
}

func (q *RabbitMQ) Push(event model.StatisticsEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		q.logger.Error("fail marshal error:" + err.Error())
		return err
	}

	err = q.channel.Publish(
		"",
		nameQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventBytes,
		},
	)
	if err != nil {
		q.logger.Error("not push to queue:" + err.Error())
		return err
	}

	return err
}

func (q *RabbitMQ) Stop() {
	err := q.channel.Close()
	if err != nil {
		q.logger.Fatal("fail to close channel:" + err.Error())
	}

	err = q.connection.Close()
	if err != nil {
		q.logger.Fatal("fail to close connection" + err.Error())
	}
}

func getUrl(config configuration.RabbitMQ) string {
	return "amqp://guest:guest@localhost:5672/"
}
