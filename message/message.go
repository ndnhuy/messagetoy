package message

type Message struct {
	// UUID is a unique identifier of message.
	// can be empty
	UUID    string
	Headers map[string]string
	Payload []byte
}
