package message

import "context"

type Publisher interface {
	Publish(topic string, messages ...*Message) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (<-chan *Message, error)
}
