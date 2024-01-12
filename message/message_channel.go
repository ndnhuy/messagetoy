package message

type MessageChannel interface {
	Send() bool
}

type SubscribableChannel interface {
	Subscribe()
}

type MessageHandler interface {
	Handle(message *Message) error
}
