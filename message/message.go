package message

import "github.com/google/uuid"

type Message struct {
	// UUID is a unique identifier of message.
	// can be empty
	UUID    string
	Headers map[string]string
	Payload []byte
}

func NewMessage(uuid string, payload []byte, headers map[string]string) *Message {
	return &Message{
		UUID:    uuid,
		Payload: payload,
		Headers: headers,
	}
}
func NewUUIDMessage(payload []byte, headers map[string]string) *Message {
	return NewMessage(uuid.New().String(), payload, headers)
}
