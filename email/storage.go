package email

import "github.com/husainaloos/notfy/status"

// Storage is an interface for storing emails
type Storage interface {
	Insert(Email, status.Info) (Email, error)
}

type inMemoryStorageData struct {
	e    Email
	info status.Info
}

// InMemoryStorage is an in-memory impelmentation of storage
type InMemoryStorage struct {
	data []inMemoryStorageData
}

// NewInMemoryStorage creates a new in-memory implementation
func NewInMemoryStorage() *InMemoryStorage { return &InMemoryStorage{} }

// Insert new email/info to a slice
func (s *InMemoryStorage) Insert(e Email, info status.Info) (Email, error) {
	e.SetID(len(s.data) + 1)
	d := inMemoryStorageData{e, info}
	s.data = append(s.data, d)
	return e, nil
}
