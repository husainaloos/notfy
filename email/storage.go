package email

type Storage interface{}

type InMemoryStorage struct{}

func NewInMemoryStorage() *InMemoryStorage { return &InMemoryStorage{} }
