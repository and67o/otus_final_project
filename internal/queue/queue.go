package rmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/and67o/otus_project/internal/interfaces"
	"time"

	"github.com/and67o/otus_project/internal/configuration"
	"github.com/and67o/otus_project/internal/model"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	logger     interfaces.Logger
	connection *amqp.Connection
	channel    *amqp.Channel
}

const nameQueue = "banner_queue"
const nameExchangeQueue = "banner_exchange_queue"

func New(config configuration.RabbitMQ, logg interfaces.Logger) (*RabbitMQ, error) {
	var err error

	res := &RabbitMQ{
		logger: logg,
	}

	_, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res.connection, err = amqp.Dial(getURL(config))
	if err != nil {
		res.logger.Fatal("fail Rabbit connection:" + err.Error())

		return nil, err
	}

	res.channel, err = res.connection.Channel()
	if err != nil {
		res.logger.Fatal("fail to open channel for Rabbit:" + err.Error())

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
		res.logger.Fatal("fail create queue for Rabbit")

		return nil, err
	}

	return res, nil
}

func (q *RabbitMQ) Push(event model.StatisticsEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		q.logger.Error("fail marshal error:" + err.Error())
		return err
	}

	err = q.channel.Publish(
		nameExchangeQueue,
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

func getURL(config configuration.RabbitMQ) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", config.User, config.Pass, config.Host, config.Port)
}
