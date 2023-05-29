package inmemory

import (
	"context"
	pubsub "github.com/dashjay/hydra-pubsub"
	"sync"
)

// inMemory, a hub in memory manages a lot of topics
type inMemory struct {
	topics map[string]*Dispatcher
	sync.Mutex
}

func New() *inMemory {
	return &inMemory{
		topics: make(map[string]*Dispatcher),
	}
}

// getDispatcher get dispatcher by topic
func (i *inMemory) getDispatcher(topic string) *Dispatcher {
	i.Lock()
	defer i.Unlock()
	d, ok := i.topics[topic]
	if !ok {
		d = NewDispatcher()
		i.topics[topic] = d
	}
	return d
}

func (i *inMemory) Publish(ctx context.Context, topic string, message interface{}) {
	d := i.getDispatcher(topic)
	m := pubsub.FromMessage(topic, message)
	d.putOneByOne(ctx, m)
}

func (i *inMemory) watch(ctx context.Context, topics []string, sub *subscription, ready chan<- struct{}) {
	defer func() {
		// first call context.Cancel()
		sub.Close()

		// then remove every topic's receiver channel from dispatcher of this topic
		for _, topic := range topics {
			d := i.getDispatcher(topic)
			d.Remove(sub.ch)
		}
		// close the channel of this subscription
		close(sub.ch)
	}()

	// add channel of this subscription to dispatcher of topic
	for _, topic := range topics {
		d := i.getDispatcher(topic)
		d.Add(sub.ch)
	}
	// ready!
	ready <- struct{}{}

	// wait until cancel()
	<-ctx.Done()
}

func (i *inMemory) Subscribe(ctx context.Context, option pubsub.SubscriptionOptions) pubsub.Subscription {
	buf := option.GetBufferSize()
	if buf < 0 {
		buf = 0
	}
	ch := make(chan pubsub.Message, buf)
	ready := make(chan struct{})
	subCtx, sub := newSubscription(ctx, ch)
	go i.watch(subCtx, option.GetAllTopics(), sub, ready)
	<-ready
	return sub
}

var _ pubsub.Pubsub = (*inMemory)(nil)
