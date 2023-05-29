package viahttp

import (
	"context"
	"encoding/json"
	"fmt"
	pubsub "github.com/dashjay/hydra-pubsub"
	"github.com/valyala/fasthttp"
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
	go fasthttp.ListenAndServe(fmt.Sprintf(":%d", port), sub.RequestHandler)
	return subCtx, sub, nil
}

func (s *subscription) Chan() <-chan pubsub.Message {
	return s.ch
}

func (s *subscription) Close() {
	s.cancel()
}

func (s *subscription) RequestHandler(ctx *fasthttp.RequestCtx) {
	var res Response
	var msg message
	err := json.Unmarshal(ctx.Request.Body(), &msg)
	if err != nil {
		res.Err = err.Error()
		bin, _ := json.Marshal(res)
		ctx.Response.SetBody(bin)
		return
	}
	s.ch <- &msg
}
