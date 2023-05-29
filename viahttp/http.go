package viahttp

import (
	"context"
	"fmt"
	pubsub "github.com/dashjay/hydra-pubsub"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	serverType = "server"
	clientType = "client"
)

// viaHTTP, a hub manages a lot of topics by http request and response
type viaHTTP struct {
	// t type of
	t string
	// port means which port this server uses
	port   int
	topics map[string]*Dispatcher
	sync.Mutex
}

func (v *viaHTTP) Publish(ctx context.Context, topic string, message interface{}) {
	if v.t == clientType {
		panic(fmt.Sprintf("you can not publish in instance type %s", v.t))
	}
	d := v.getDispatcher(topic)
	m := pubsub.FromMessage(topic, message)
	d.putOneByOne(ctx, m)
}

// getDispatcher get dispatcher by topic
func (v *viaHTTP) getDispatcher(topic string) *Dispatcher {
	v.Lock()
	defer v.Unlock()
	d, ok := v.topics[topic]
	if !ok {
		d = NewDispatcher()
		v.topics[topic] = d
	}
	return d
}

func (v *viaHTTP) Subscribe(ctx context.Context, option pubsub.SubscriptionOptions) pubsub.Subscription {
	if v.t == serverType {
		panic(fmt.Sprintf("you can not subscribe in instance type %s", v.t))
	}
	buf := option.GetBufferSize()
	if buf < 0 {
		buf = 0
	}
	ch := make(chan pubsub.Message, buf)
	ready := make(chan struct{})
	subCtx, sub, err := newSubscription(ctx, option.(*SubscriptionOptions).Port, ch)
	if err != nil {
		logrus.WithError(err).Errorf("create new Subscription error")
		return nil
	}
	go v.watch(subCtx, option.GetAllTopics(), sub, ready)
	<-ready
	return sub
}

func (v *viaHTTP) watch(ctx context.Context, topics []string, sub *subscription, ready chan struct{}) {
	defer func() {
		// first call context.Cancel()
		sub.Close()

		// then remove every topic's receiver channel from dispatcher of this topic
		for _, topic := range topics {
			d := v.getDispatcher(topic)
			d.Remove(sub.clientAddr)
		}
		// close the channel of this subscription
		close(sub.ch)
	}()

	// add channel of this subscription to dispatcher of topic
	for _, topic := range topics {
		d := v.getDispatcher(topic)
		d.Add(sub.clientAddr)
	}
	// ready!
	ready <- struct{}{}

	// wait until cancel()
	<-ctx.Done()
}

func NewServer(port int) *viaHTTP {
	return &viaHTTP{
		port:   port,
		topics: make(map[string]*Dispatcher),
	}
}

func NewClient(serverAddr string) *viaHTTP {
	return &viaHTTP{
		port: port,
		t:    "client",
	}
}

var _ pubsub.Pubsub = (*viaHTTP)(nil)
