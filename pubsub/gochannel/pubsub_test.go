package gochannel

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ndnhuy/messagetoy/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublishSubscribe_not_persistent(t *testing.T) {
	// create pub sub
	pub, sub := createPubSub()

	// subscribe
	msgs, err := sub.Subscribe(context.TODO(), "topic_test")
	require.NoError(t, err)

	// publish
	msgToPublish := &message.Message{
		Headers: map[string]string{
			"id": "1",
		},
		Payload: []byte("hello"),
	}
	err = pub.Publish("topic_test", msgToPublish)
	require.NoError(t, err)

	// publish test messages
	publishedMsg, err := publishTestMessages(1, pub, "topic_test")
	require.NoError(t, err)

	// verify if the messages received from the subscriber equals to the messages that was published
	receivedMsgs := []*message.Message{}
	select {
	case receivedMsg, ok := <-msgs:
		if ok {
			// assert.Equal(t, msgToPublish.Payload, receivedMsg.Payload)
		} else {
			// assert.Fail(t, "subscriber fail to receive message")
		}
	case <-time.After(3 * time.Second):
		assert.Fail(t, "subscribe timeout, no message delivered to subscriber")
	}
}

func createPubSub() (message.Publisher, message.Subscriber) {
	gochannel := NewGoChannel()
	return gochannel, gochannel
}

// publish simple messages without payload
func publishTestMessages(msgCount int, pub message.Publisher, topic string) ([]*message.Message, error) {
	publishedMsg := []*message.Message{}
	for i := 0; i < msgCount; i++ {
		msg := &message.Message{
			UUID: uuid.New().String(),
		}
		err := pub.Publish(topic, msg)
		if err != nil {
			return nil, err
		}
		publishedMsg = append(publishedMsg, msg)
	}
	return publishedMsg, nil
}
