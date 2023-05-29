package inmemory

import (
	"context"
	"sync"

	pubsub "github.com/dashjay/hydra-pubsub"
)

// Dispatcher is in charge of a topic
type Dispatcher struct {
	channels []chan<- pubsub.Message
	sync.RWMutex
}

// NewDispatcher create new Dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{channels: make([]chan<- pubsub.Message, 0, 0)}
}

// Add one subscriber to the topic
func (d *Dispatcher) Add(ch chan<- pubsub.Message) {
	d.Lock()
	defer d.Unlock()
	d.channels = append(d.channels, ch)
}

// Remove removes a channel from channels
func (d *Dispatcher) Remove(ch chan pubsub.Message) {
	d.Lock()
	defer d.Unlock()
	nch := make([]chan<- pubsub.Message, len(d.channels)-1)
	for i := 0; i < len(d.channels); i++ {
		if ch != d.channels[i] {
			nch = append(nch, d.channels[i])
		}
	}
	d.channels = nch
}

// putOneByOne like it's name, put message one to subscriber one by one
func (d *Dispatcher) putOneByOne(ctx context.Context, message pubsub.Message) {
	d.RLock()
	channels := d.channels
	d.RUnlock()
	for _, ch := range channels {
		select {
		case <-ctx.Done():
			return
		case ch <- message:
		}
	}
}
