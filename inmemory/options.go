package inmemory

import (
	pubsub "github.com/dashjay/hydra-pubsub"
)

type SubscriptionOptions struct {
	BufferSize int
	Topics     []string
}

func (s *SubscriptionOptions) GetBufferSize() int {
	return s.BufferSize
}

func (s *SubscriptionOptions) GetAllTopics() []string {
	return s.Topics
}

var _ pubsub.SubscriptionOptions = (*SubscriptionOptions)(nil)

func NewSubscriptionOptions(bufSize int, topics ...string) *SubscriptionOptions {
	return &SubscriptionOptions{BufferSize: bufSize, Topics: topics}
}
