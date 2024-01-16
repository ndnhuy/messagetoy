package gochannel

import (
	"context"

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
	s := newSubscriber(ctx)
	if _, ok := g.subscribers[topic]; !ok {
		g.subscribers[topic] = make([]*subscriber, 0)
	}
	g.subscribers[topic] = append(g.subscribers[topic], s)
	return s.msgQueue, nil
}

func (g *GoChannel) Close() error {
	for _, subs := range g.subscribers {
		for _, sub := range subs {
			sub.close()
		}
	}
	return nil
}
