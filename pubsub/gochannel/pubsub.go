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
	// subscribers by topic
	subscribers     map[string][]*subscriber
	subscribersLock sync.RWMutex

	// lock that guards the closed variable,
	// make sure no subscriber is created while channel is closing
	closedLock sync.RWMutex
	closed     bool
}

func NewGoChannel() *GoChannel {
	return &GoChannel{
		subscribers: make(map[string][]*subscriber),
	}
}

// send message to all subscribers
// return
// - err if channel is already closed
// TODO implement ack by subscriber
func (g *GoChannel) Publish(topic string, messages ...*message.Message) (err error) {
	if g.IsClosed() {
		return errClosedChannel
	}
	subs := g.subscribers[topic]
	for _, msg := range messages {
		for _, sub := range subs {
			go sub.send(msg)
		}
	}

	return err
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

func (g *GoChannel) IsClosed() bool {
	g.closedLock.RLock()
	defer g.closedLock.RUnlock()
	return g.closed
}
