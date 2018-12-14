package messaging

// InMemoryBroker is a broker that runs in memory
type InMemoryBroker struct{ C chan []byte }

// NewInMemoryBroker creates new instance of InMemoryBroker
func NewInMemoryBroker() *InMemoryBroker {
	return &InMemoryBroker{make(chan []byte, 100000)}
}

// Publish to the in-memory channel
func (b *InMemoryBroker) Publish(bb []byte) error {
	b.C <- bb
	return nil
}

// Close the in-memory channel
func (b *InMemoryBroker) Close() error {
	close(b.C)
	return nil
}

// Consume messages from broker
func (b *InMemoryBroker) Consume() ([]byte, error) {
	return <-b.C, nil
}
