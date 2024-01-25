package gochannel

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
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

	// publish test messages
	msgCnt := 100
	publishedMsgs, err := publishTestMessages(msgCnt, pub, "topic_test")
	require.NoError(t, err)

	receivedMsgs, _ := bulkRead(msgs, msgCnt, 3*time.Second)
	// assert that sent messages and received messages are same
	assertAllMsgsAreReceived(t, publishedMsgs, receivedMsgs)

	// close pub sub
	pub.Close()
	sub.Close()

	// assert pub/sub is closed
	select {
	case _, open := <-msgs:
		assert.False(t, open)
	default:
		t.Error("messages channel is not closed")
	}
}

func TestSubscriber_race_condition_when_closing(t *testing.T) {
	for i := 0; i < 1000; i++ {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			_, sub := createPubSub()
			go func() {
				err := sub.Close()
				require.NoError(t, err)
			}()
			msgQueue, err := sub.Subscribe(context.TODO(), "test")
			if errors.Is(err, errClosedChannel) {
				fmt.Println(errClosedChannel.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, msgQueue)
				select {
				case _, open := <-msgQueue:
					assert.False(t, open)
				case <-time.After(3 * time.Second):
					assert.Fail(t, "timeout, the subscriber should be closed by now")
				}
			}
		})
	}
}

func TestPublish_race_condition_when_closing(t *testing.T) {
	testCnt := 10
	for i := 0; i < testCnt; i++ {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			pub, _ := createPubSub()
			go func() {
				err := pub.Publish("test", message.NewUUIDMessage([]byte("hello"), nil))
				fmt.Println("fail to publish: " + err.Error())
			}()

			err := pub.Close()
			require.NoError(t, err)
		})
	}
}

func assertAllMsgsAreReceived(t *testing.T, publishedMsgs []*message.Message, receivedMsgs []*message.Message) {
	assert.Equal(t, len(publishedMsgs), len(receivedMsgs))
	sentIDs := toMsgIDs(publishedMsgs)
	receivedIDs := toMsgIDs(receivedMsgs)
	sort.Strings(sentIDs)
	sort.Strings(receivedIDs)
	assert.Equal(t, sentIDs, receivedIDs)
}

func bulkRead(msgs <-chan *message.Message, limit int, timeout time.Duration) (receivedMsgs []*message.Message, allRead bool) {
ReceiveMsgLoop:
	for len(receivedMsgs) < limit {
		select {
		case receivedMsg, ok := <-msgs:
			if ok {
				receivedMsgs = append(receivedMsgs, receivedMsg)
			} else {
				// channel is closed, just break
				allRead = false
				break ReceiveMsgLoop
			}
		case <-time.After(timeout):
			break ReceiveMsgLoop
		}
	}
	return receivedMsgs, limit == len(receivedMsgs)
}

func toMsgIDs(msgs []*message.Message) []string {
	ids := []string{}
	for _, msg := range msgs {
		ids = append(ids, msg.UUID)
	}
	return ids
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
			UUID:    uuid.New().String(),
			Payload: []byte(strconv.Itoa(i)),
		}
		err := pub.Publish(topic, msg)
		if err != nil {
			return nil, err
		}
		publishedMsg = append(publishedMsg, msg)
	}
	return publishedMsg, nil
}
