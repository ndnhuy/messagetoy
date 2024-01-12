package gochannel

import (
	"context"
	"fmt"

	"github.com/ndnhuy/messagetoy/message"
)

type GoChannel struct {
	subscribers map[string][]*subscriber // subscribers by topic
}

func NewGoChannel() *GoChannel {
	return &GoChannel{
		subscribers: make(map[string][]*subscriber),
	}
}

func (g *GoChannel) Publish(topic string, messages ...*message.Message) error {
	subs := g.subscribers[topic]
	for _, msg := range messages {
		for _, sub := range subs {
			go sub.send(msg)
		}
	}
	return nil
}
func (g *GoChannel) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	s := &subscriber{
		ctx:      ctx,
		msgQueue: make(chan *message.Message),
	}
	if _, ok := g.subscribers[topic]; !ok {
		g.subscribers[topic] = make([]*subscriber, 0)
	}
	g.subscribers[topic] = append(g.subscribers[topic], s)
	return s.msgQueue, nil
}

type subscriber struct {
	ctx      context.Context
	msgQueue chan *message.Message
}

func (s *subscriber) send(msg *message.Message) {
	s.msgQueue <- msg
	fmt.Println("sent msg to subscriber: " + string(msg.Payload))
}
