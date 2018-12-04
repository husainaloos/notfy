package email

import (
	"context"
	"errors"
)

var (
	errStorageNotFound = errors.New("item not found")
)

type Storage interface {
	insert(context.Context, Email) (Email, error)
}

type MemoryStorage struct {
	emails []Email
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		emails: make([]Email, 0),
	}
}

func (s *MemoryStorage) insert(ctx context.Context, e Email) (Email, error) {
	e.SetID(len(s.emails) + 1)
	s.emails = append(s.emails, e)
	return e, nil
}
