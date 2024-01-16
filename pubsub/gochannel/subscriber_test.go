package gochannel

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ndnhuy/messagetoy/message"
)

func TestSubscriber_close_subscriber_in_middle_of_sending_msg(t *testing.T) {
	msgCnt := 1
	sub := newBufferSubscriber(context.TODO(), msgCnt)
	msg := message.NewMessage([]byte("test"), nil)
	sub.send(msg)
	sub.close()

	select {
	case receiveMsg := <-sub.msgQueue:
		fmt.Println("receive msg " + string(receiveMsg.Payload))
	case <-time.After(3 * time.Second):
	}

	// var wg sync.WaitGroup
	// wg.Add(msgCnt)
	// go func() {
	// 	for i := 0; i < msgCnt; i++ {
	// 		msg := message.NewMessage([]byte("offset: "+strconv.Itoa(i)), nil)
	// 		time.Sleep(100 * time.Millisecond)
	// 		go func() {
	// 			defer wg.Done()
	// 			sub.send(msg)
	// 		}()
	// 	}
	// }()

	// time.Sleep(500 * time.Millisecond)
	// go sub.close()

	// wg.Wait()

	//	i := 0
	//
	// ReceiveMsgLoop:
	//
	//	for i < msgCnt {
	//		select {
	//		case receiveMsg := <-sub.msgQueue:
	//			fmt.Println("receive msg " + string(receiveMsg.Payload))
	//			i++
	//		case <-time.After(3 * time.Second):
	//			break ReceiveMsgLoop
	//		}
	//	}
}
