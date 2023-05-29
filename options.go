package pubsub

// SubscriptionOptions can help subscriber do subscribe on more than one topic
// and use buffer size as a parameter
type SubscriptionOptions interface {
	GetBufferSize() int
	GetAllTopics() []string
}
