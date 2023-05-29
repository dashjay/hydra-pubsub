package client

import (
	"context"
	pubsub "github.com/dashjay/hydra-pubsub"
	"net"
)

// subscription will be hold by subscriber to get the published message
type subscription struct {
	port int

	ch         chan pubsub.Message
	clientAddr string

	cancel context.CancelFunc
}

var _ pubsub.Subscription = (*subscription)(nil)

// newSubscription new a subscription
func newSubscription(ctx context.Context, port int, ch chan pubsub.Message) (context.Context, *subscription, error) {
	subCtx, cancel := context.WithCancel(ctx)
	sub := &subscription{cancel: cancel, ch: ch}
	go func() {
		net.Dial("ip4","")
	}()
	return subCtx, sub, nil
}

func (s *subscription) Chan() <-chan pubsub.Message {
	return s.ch
}

func (s *subscription) Close() {
	s.cancel()
}
