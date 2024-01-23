package gochannel

import (
	"context"
	"errors"
	"fmt"

	"github.com/ndnhuy/messagetoy/message"
)

type subscriber struct {
	ctx      context.Context
	msgQueue chan *message.Message
	closed   bool
}

func newSubscriber(ctx context.Context) *subscriber {
	return &subscriber{
		ctx:      ctx,
		msgQueue: make(chan *message.Message),
	}
}

func newBufferSubscriber(ctx context.Context, buffer int) *subscriber {
	return &subscriber{
		ctx:      ctx,
		msgQueue: make(chan *message.Message, buffer),
	}
}

func (s *subscriber) send(msg *message.Message) error {
	if s.closed {
		return errors.New("subscriber is closed")
	}
	s.msgQueue <- msg
	fmt.Printf("sent msg to subscriber: {uuid: %s, payload: %s}\n", msg.UUID, string(msg.Payload))
	return nil
}

func (s *subscriber) close() {
	if !s.closed {
		close(s.msgQueue)
		s.closed = true
	}
}
