package rabbitmq

import (
	"context"
	"errors"
	"fmt"

	"github.com/ndnhuy/messagetoy/message"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscriber struct {
	conn *amqp.Connection
}

func NewAMQPSubscriber() (*Subscriber, error) {
	// create new connection to rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, errors.New("fail to connect to rabbitmq")
	}
	return &Subscriber{
		conn: conn,
	}, nil
}

func (sub *Subscriber) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	// open new channel
	ch, err := sub.conn.Channel()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	q, err := ch.QueueDeclare(
		topic,
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		"",
		topic,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.New("fail to bind the consumer queue to exchange: " + err.Error())
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.New("cannot consume from channel")
	}

	// consume and process message
	out := make(chan *message.Message)
	go func() {
		defer func() {
			fmt.Println("stop consuming from rabbitmq channel")
			close(out)
			if err := ch.Close(); err != nil {
				fmt.Println("failed to close channel")
			}
		}()

		for amqpMsg := range msgs {
			headers := map[string]string{}
			for k, v := range amqpMsg.Headers {
				headers[k] = v.(string)
			}
			msg := message.NewMessage(headers["uuid"], amqpMsg.Body, headers)
			out <- msg
			fmt.Println("message sent to consumer")
		}
	}()

	return out, nil
}

func (sub *Subscriber) Close() error {
	return sub.conn.Close()
}
