package email

import (
	"context"
)

type Storage interface {
	insert(context.Context, Email) (Email, error)
	get(context.Context, int) (Email, bool, error)
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

func (s *MemoryStorage) get(ctx context.Context, id int) (Email, bool, error) {
	for _, v := range s.emails {
		if v.ID() == id {
			return v, true, nil
		}
	}
	return Email{}, false, nil
}
