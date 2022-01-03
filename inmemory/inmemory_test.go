package inmemory

import (
	"context"
	pubsub "github.com/dashjay/hydra-pubsub"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestInMemory(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	inm := New()
	sub := inm.Subscribe(ctx, pubsub.NewSubscriptionOptions(64, "love", "hate"))
	ch := sub.Chan()
	inm.Publish(ctx, "love", 1)
	msg := <-ch
	require.Equal(t, 1, msg.Message())
	inm.Publish(ctx, "love", 2)
	msg = <-ch
	require.Equal(t, 2, msg.Message())
	cancel()
	msg, more := <-ch
	require.Equal(t, false, more)
}

func TestInMemoryMultiSubscriber(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	inm := New()
	sub1 := inm.Subscribe(ctx, pubsub.NewSubscriptionOptions(64, "love", "hate"))
	sub2 := inm.Subscribe(ctx, pubsub.NewSubscriptionOptions(64, "love", "hate"))
	inm.Publish(ctx, "love", 666)
	m1 := <-sub1.Chan()
	t.Logf("sub1 receive message %v", m1.Message())
	m2 := <-sub2.Chan()
	t.Logf("sub2 receive message %v", m2.Message())
}

func BenchmarkInMemoryPublishMany(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	inm := New()
	sub1, sub2 :=
		inm.Subscribe(ctx, pubsub.NewSubscriptionOptions(64, "love", "hate")),
		inm.Subscribe(ctx, pubsub.NewSubscriptionOptions(64, "love", "hate"))

	done := make(chan struct{})
	go func() {
		for i := 0; i < b.N; i++ {
			inm.Publish(ctx, "love", 666)
		}
		done <- struct{}{}
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		ch := sub1.Chan()
		for i := 0; i < b.N; i++ {
			<-ch
		}
		sub1.Close()
	}()
	go func() {
		defer wg.Done()
		ch := sub2.Chan()
		for i := 0; i < b.N; i++ {
			<-ch
		}
		sub2.Close()
	}()
	<-done
	wg.Wait()
}

func BenchmarkInMemorySubscriptionMany(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	inm := New()
	opt := pubsub.NewSubscriptionOptions(64, "love", "hate")

	var subscriptions = make([]pubsub.Subscription, b.N)
	for i := 0; i < b.N; i++ {
		subscriptions[i] = inm.Subscribe(ctx, opt)
	}
	inm.Publish(ctx, "love", 666)
	for i := 0; i < b.N; i++ {
		<-subscriptions[i].Chan()
	}
}
