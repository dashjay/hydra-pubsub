package server

import (
	"context"
	"encoding/json"
	pubsub "github.com/dashjay/hydra-pubsub"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
)

// Dispatcher is in charge of a topic
type Dispatcher struct {
	clients []net.Conn
	sync.RWMutex
}

// NewDispatcher create new Dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

// Remove removes a channel from clients
func (d *Dispatcher) Remove(conn net.Conn) {
	d.Lock()
	nch := make([]net.Conn, len(d.clients)-1)
	for i := 0; i < len(d.clients); i++ {
		if conn != d.clients[i] {
			nch = append(nch, d.clients[i])
		}
	}
	d.clients = nch
	d.Unlock()
}

// Add add a server addr into server slice
func (d *Dispatcher) Add(conn net.Conn) {
	d.Lock()
	d.clients = append(d.clients, conn)
	d.Unlock()
}

// putOneByOne like it's name, put message one to subscriber one by one
func (d *Dispatcher) putOneByOne(ctx context.Context, message pubsub.Message) {
	d.RLock()
	clients := d.clients
	d.RUnlock()

	bin, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	for idx := range clients {
		ch := clients[idx]
		select {
		case <-ctx.Done():
			return
		default:
			try := func() error {
				_, err := ch.Write(bin)
				return err
			}
			for i := 0; i < 10; i++ {
				err := try()
				if err != nil {
					logrus.WithError(err).Errorf("write to client [%s] error", ch.RemoteAddr())
					continue
				}
				break
			}
		}
	}
}
