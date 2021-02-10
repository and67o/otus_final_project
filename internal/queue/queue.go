package rmq

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/and67o/otus_project/internal/configuration"
	"github.com/and67o/otus_project/internal/interfaces"
	"github.com/and67o/otus_project/internal/model"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	configuration configuration.RabbitMQ
	connection    *amqp.Connection
	channel       *amqp.Channel
}

func New(config configuration.RabbitMQ) (interfaces.Queue, error) {
	var r RabbitMQ

	connection, err := amqp.Dial(getURL(config))
	if err != nil {
		return nil, errors.New("connection error:" + err.Error())
	}

	r.configuration = config
	r.connection = connection

	return &r, nil
}

func (q *RabbitMQ) Push(event model.StatisticsEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return errors.New("fail marshal error:" + err.Error())
	}

	err = q.channel.Publish(
		q.configuration.Name,
		q.configuration.Key,
		false,
		false,
		amqp.Publishing{
			ContentType:     "application/json",
			Body:            eventBytes,
			ContentEncoding: "utf8",
		},
	)
	if err != nil {
		return errors.New("not push to queue:" + err.Error())
	}

	return nil
}

func (q *RabbitMQ) CloseConnection() error {
	return q.connection.Close()
}

func (q *RabbitMQ) OpenChanel() error {
	channel, err := q.connection.Channel()
	if err != nil {
		return fmt.Errorf("fail to open channel for Rabbit: %w", err)
	}

	err = channel.ExchangeDeclare(
		q.configuration.Name,
		q.configuration.Kind,
		q.configuration.Durable,
		q.configuration.AutoDelete,
		q.configuration.Internal,
		q.configuration.NoWait,
		nil,
	)

	if err != nil {
		return fmt.Errorf("fail create queue for Rabbit: %w", err)
	}
	q.channel = channel

	return nil
}

func (q *RabbitMQ) CloseChannel() error {
	return q.channel.Close()
}

func getURL(config configuration.RabbitMQ) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", config.User, config.Pass, config.Host, config.Port)
}
