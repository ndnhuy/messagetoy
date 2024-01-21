package gochannel

import (
	"context"
	"errors"
	"sync"

	"github.com/ndnhuy/messagetoy/message"
)

var (
	errClosedChannel = errors.New("channel is closed")
)

type GoChannel struct {
	subscribers     map[string][]*subscriber // subscribers by topic
	subscribersLock sync.RWMutex

	closedLock sync.RWMutex
	closed     bool
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
	g.closedLock.RLock()
	defer g.closedLock.RUnlock()
	if g.closed {
		return nil, errClosedChannel
	}
	s := newSubscriber(ctx)
	g.subscribersLock.Lock()
	defer g.subscribersLock.Unlock()
	return g.addSubscriber(topic, s)
}

func (g *GoChannel) addSubscriber(topic string, sub *subscriber) (<-chan *message.Message, error) {
	if _, ok := g.subscribers[topic]; !ok {
		g.subscribers[topic] = make([]*subscriber, 0)
	}
	g.subscribers[topic] = append(g.subscribers[topic], sub)
	return sub.msgQueue, nil
}

func (g *GoChannel) Close() error {
	g.closedLock.Lock()
	defer g.closedLock.Unlock()
	for _, subs := range g.subscribers {
		for _, sub := range subs {
			sub.close()
		}
	}
	g.closed = true
	return nil
}
