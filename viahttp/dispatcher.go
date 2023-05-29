package viahttp

import (
	"context"
	"encoding/json"
	pubsub "github.com/dashjay/hydra-pubsub"
	"github.com/valyala/fasthttp"
	"sync"
)

// Dispatcher is in charge of a topic
type Dispatcher struct {
	clients []string
	sync.RWMutex
}

// NewDispatcher create new Dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

// Remove removes a channel from clients
func (d *Dispatcher) Remove(addr string) {
	d.Lock()
	nch := make([]string, len(d.clients)-1)
	for i := 0; i < len(d.clients); i++ {
		if addr != d.clients[i] {
			nch = append(nch, d.clients[i])
		}
	}
	d.clients = nch
	d.Unlock()
}

// Add add a server addr into server slice
func (d *Dispatcher) Add(addr string) {
	d.Lock()
	d.clients = append(d.clients, addr)
	d.Unlock()
}

// putOneByOne like it's name, put message one to subscriber one by one
func (d *Dispatcher) putOneByOne(ctx context.Context, message pubsub.Message) []Response {
	d.RLock()
	clients := d.clients
	d.RUnlock()

	var responses []Response
	for idx := range clients {
		ch := clients[idx]
		select {
		case <-ctx.Done():
			return responses
		default:
			req := fasthttp.AcquireRequest()
			bin, err := json.Marshal(message)
			if err != nil {
				responses = append(responses, Response{Err: err.Error()})
				// never use defer in loop
				fasthttp.ReleaseRequest(req)
				continue
			}
			req.SetBody(bin)
			req.SetRequestURI(ch)
			resp := fasthttp.AcquireResponse()
			err = fasthttp.Do(req, resp)
			if err != nil {
				responses = append(responses, Response{Err: err.Error()})
				fasthttp.ReleaseRequest(req)
				fasthttp.ReleaseResponse(resp)
				continue
			}
			var out Response
			err = json.Unmarshal(resp.Body(), &out)
			if err != nil {
				responses = append(responses, Response{Err: err.Error()})
				fasthttp.ReleaseRequest(req)
				fasthttp.ReleaseResponse(resp)
				continue
			}
			if out.Err != "" {
				responses = append(responses, Response{Err: out.Err})
				fasthttp.ReleaseRequest(req)
				fasthttp.ReleaseResponse(resp)
				continue
			}
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
		}
	}
	return responses
}
