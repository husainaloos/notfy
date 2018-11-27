package messaging

// Publisher is a publisher to the message broker of some sort
type Publisher interface {
	Publish([]byte) error
	Close() error
}

// NilPublisher is a publisher that does nothing
type NilPublisher struct{}

// Publish does nothing
func (n NilPublisher) Publish([]byte) error { return nil }

// Close does nothing
func (n NilPublisher) Close() error { return nil }

// InMemoryPublisher is an in-memory channel that implements the publisher interface
type InMemoryPublisher struct{ C chan []byte }

// NewInMemoryPublisher creates a new in-memory publisher
func NewInMemoryPublisher() *InMemoryPublisher {
	return &InMemoryPublisher{make(chan []byte, 100000)}
}

// Publish to the in-memory channel
func (b *InMemoryPublisher) Publish(bb []byte) error {
	b.C <- bb
	return nil
}

// Close the in-memory channel
func (b *InMemoryPublisher) Close() error {
	close(b.C)
	return nil
}

// ErrPublisher is an implementation of the publisher that returns an error if the error is set
type ErrPublisher struct{ Err error }

// Publish returns an error if an error is set
func (p ErrPublisher) Publish([]byte) error { return p.Err }

// Close returns nil
func (p ErrPublisher) Close() error { return nil }
