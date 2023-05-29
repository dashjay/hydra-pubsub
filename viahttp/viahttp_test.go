package viahttp

import (
	"context"
	"testing"
)

func TestViaHttp(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server := NewServer(8080)

	client := NewClient()
	inm.Publish(ctx, "love", "fuck")
}
