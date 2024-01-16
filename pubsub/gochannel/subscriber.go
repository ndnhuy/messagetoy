package gochannel

import (
	"context"
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

func (s *subscriber) send(msg *message.Message) {
	s.msgQueue <- msg
	fmt.Printf("sent msg to subscriber: {uuid: %s, payload: %s}\n", msg.UUID, string(msg.Payload))
}

func (s *subscriber) close() {
	if !s.closed {
		close(s.msgQueue)
		s.closed = true
	}
}
