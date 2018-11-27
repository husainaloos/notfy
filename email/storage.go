package email

import (
	"errors"

	"github.com/husainaloos/notfy/status"
)

var (
	errStorageNotFound = errors.New("id not found")
)

type statusEmail struct {
	Email
	statusID int
}

// Storage is an interface for storing emails
type Storage interface {
	insert(Email, status.Info) (Email, error)
	get(int) (statusEmail, error)
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

func (s *InMemoryStorage) insert(e Email, info status.Info) (Email, error) {
	e.SetID(len(s.data) + 1)
	d := inMemoryStorageData{e, info}
	s.data = append(s.data, d)
	return e, nil
}

func (s *InMemoryStorage) get(id int) (statusEmail, error) {
	for _, v := range s.data {
		if v.e.ID() == id {
			return statusEmail{v.e, v.info.ID()}, nil
		}
	}
	return statusEmail{}, errStorageNotFound
}
