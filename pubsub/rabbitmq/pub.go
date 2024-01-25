package rabbitmq

import (
	"context"
	"errors"

	"github.com/ndnhuy/messagetoy/message"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn *amqp.Connection
}

func NewAMQPPublisher() (*Publisher, error) {
	// create new connection to rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, errors.New("fail to connect to rabbitmq")
	}
	return &Publisher{
		conn: conn,
	}, nil
}

func (pub *Publisher) Publish(topic string, messages ...*message.Message) error {
	// open new channel
	ch, err := pub.conn.Channel()
	if err != nil {
		return err
	}
	err = ch.ExchangeDeclare(
		topic,
		"fanout",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for _, msg := range messages {
		headers := make(amqp.Table, len(msg.Headers)+1) // +1 len for uuid
		for k, v := range msg.Headers {
			headers[k] = v
		}
		headers["uuid"] = msg.UUID
		err := ch.PublishWithContext(context.Background(), topic, "", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg.Payload),
			Headers:     headers,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (pub *Publisher) Close() error {
	return pub.conn.Close()
}
