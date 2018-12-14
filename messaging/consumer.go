package messaging

// Consumer is an interface for anything that can consumer messages
type Consumer interface {
	Consume() ([]byte, error)
}

// SubscribeFunc is method invoked when subscribing
type SubscribeFunc func([]byte)

// Subscriber is interface for anything that one can subscribe to
type Subscriber interface {

	//Subscribe to the subscriber
	Subscribe(SubscribeFunc) error
}
