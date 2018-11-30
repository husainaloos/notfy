package messaging

// Consumer is an interface for anything that can consumer messages
type Consumer interface {
	Consume() ([]byte, error)
}
