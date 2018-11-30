package messaging

type InMemoryBroker struct{ C chan []byte }

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

func (b *InMemoryBroker) Consume() ([]byte, error) {
	return <-b.C, nil
}
