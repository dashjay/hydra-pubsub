package pubsub

import (
	"context"
)

// Pubsub is a hub use to publish a series of message, in different of medium like memory, tcp, or databases
type Pubsub interface {
	// Publish a message so that the hub will push these message one by one to subscriber.
	Publish(ctx context.Context, topic string, message interface{})

	// Subscribe input a topic and get a handle(Subscription) to receive message
	Subscribe(ctx context.Context, option SubscriptionOptions) Subscription
}

// Subscription is return by Subscribe, like a handle to receive the message or close it.
type Subscription interface {
	// Close should implement to wait all message send over
	Close()

	// Chan returns a <-chan like a handle to read message
	Chan() <-chan Message
}

// Message interface abstract what the msg do
type Message interface {
	// Record which topic this message recieved from
	Topic() string

	// Message can get the message itself
	Message() interface{}

	// Error can be used to pass an error back to subscription
	Error() error

	// Unmarshal can use a predefined certain way to unmarshal message into a value
	Unmarshal(v interface{}) error
}
