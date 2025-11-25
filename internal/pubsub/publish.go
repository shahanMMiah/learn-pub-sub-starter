package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	DURABLE   SimpleQueueType = "durable"
	TRANSIENT SimpleQueueType = "transient"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {

	jsondata, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ContentType: "application/json", Body: jsondata})
	if err != nil {
		return err
	}

	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	isDurable := false
	autoDel := true
	isExclusive := true

	if queueType == DURABLE {
		isDurable = true
		autoDel = false
		isExclusive = false
	}

	queue, err := channel.QueueDeclare(queueName, isDurable, autoDel, isExclusive, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = channel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return channel, queue, nil
}
