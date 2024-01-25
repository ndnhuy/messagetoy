package rabbitmq

import (
	"context"
	"errors"

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
	return nil, nil
}

func (sub *Subscriber) Close() error {
	return nil
}
