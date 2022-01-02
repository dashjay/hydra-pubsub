package inmemory

import (
	"context"
	pubsub "github.com/dashjay/hydra-pubsub"
)

// subscription will be hold by subscriber to get the published message
type subscription struct {
	ch     chan pubsub.Message
	cancel context.CancelFunc
}

var _ pubsub.Subscription = (*subscription)(nil)

// newSubscription new a subscription
func newSubscription(ctx context.Context, ch chan pubsub.Message) (context.Context, *subscription) {
	subCtx, cancel := context.WithCancel(ctx)
	return subCtx, &subscription{ch: ch, cancel: cancel}
}

func (s *subscription) Chan() <-chan pubsub.Message {
	return s.ch
}

func (s *subscription) Close() {
	s.cancel()
}
