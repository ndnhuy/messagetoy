package rabbitmq

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ndnhuy/messagetoy/message"
	"github.com/stretchr/testify/require"
)

// Publish test messages to broker, then Subscribe all of them and assert if received messages are the messages which has been sent
//
// Before running the test, we need to start the server of rabbitmq by running this command
// docker-compose up -d
//
// We can access to rabbitmq docker to inspect the server as following:
//
// sh to rabbitmq: docker exec -it rabbitmq sh
//
// enable rabbitmqadmin: rabbitmq-plugins enable rabbitmq_management
//
// list exchanges: rabbitmqadmin list exchanges
//
// list queues: rabbitmqadmin list queues
//
// get messages from the queue: rabbitmqadmin get queue=topic_test_queue count=9
//
// declare queue: rabbitmqadmin declare queue name=topic_test_queue
// list bindings: rabbitmqadmin list bindings
//
// bind queue: rabbitmqadmin declare binding source="topic_test" destination_type="queue" destination="topic_test_queue" routing_key=""
func TestPublishSubscribe_not_persistent(t *testing.T) {
	// create pub sub
	pub, sub, err := createPubSub()
	require.NoError(t, err)

	// subscribe
	msgs, err := sub.Subscribe(context.TODO(), "topic_test")
	require.NoError(t, err)

	// publish test messages
	msgCnt := 10
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

func createPubSub() (message.Publisher, message.Subscriber, error) {
	pub, err := NewAMQPPublisher()
	if err != nil {
		return nil, nil, err
	}
	sub, err := NewAMQPSubscriber()
	if err != nil {
		return nil, nil, err
	}
	return pub, sub, err
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

func assertAllMsgsAreReceived(t *testing.T, publishedMsgs []*message.Message, receivedMsgs []*message.Message) {
	assert.Equal(t, len(publishedMsgs), len(receivedMsgs))
	sentIDs := toMsgIDs(publishedMsgs)
	receivedIDs := toMsgIDs(receivedMsgs)
	sort.Strings(sentIDs)
	sort.Strings(receivedIDs)
	assert.Equal(t, sentIDs, receivedIDs)
}

func toMsgIDs(msgs []*message.Message) []string {
	ids := []string{}
	for _, msg := range msgs {
		ids = append(ids, msg.UUID)
	}
	return ids
}
