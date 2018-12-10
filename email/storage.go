package email

import (
	"context"
	"fmt"
)

type Storage interface {
	insert(context.Context, Email) (Email, error)
	get(context.Context, int) (Email, bool, error)
	update(context.Context, Email) (Email, bool, error)
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

func (s *MemoryStorage) remove(ctx context.Context, e Email) (Email, bool, error) {
	for i, v := range s.emails {
		if v.ID() == e.ID() {
			s.emails = append(s.emails[:i], s.emails[i+1:]...)
			return v, true, nil
		}
	}

	return Email{}, false, nil
}

func (s *MemoryStorage) get(ctx context.Context, id int) (Email, bool, error) {
	for _, v := range s.emails {
		if v.ID() == id {
			return v, true, nil
		}
	}
	return Email{}, false, nil
}

func (s *MemoryStorage) update(ctx context.Context, e Email) (Email, bool, error) {
	_, ok, err := s.get(ctx, e.ID())
	if err != nil || !ok {
		return Email{}, ok, err
	}

	id := e.ID()
	from := e.From()
	to := []string{}
	cc := []string{}
	bcc := []string{}
	subject := e.Subject()
	body := e.Body()

	for _, r := range e.To() {
		to = append(to, r.String())
	}
	for _, r := range e.CC() {
		cc = append(cc, r.String())
	}
	for _, r := range e.BCC() {
		bcc = append(bcc, r.String())
	}

	newE, err := New(id, from.String(), to, cc, bcc, subject, body)
	if err != nil {
		return Email{}, true, fmt.Errorf("failed to create new email: %v", err)
	}

	for _, v := range e.StatusHistory() {
		newE.AddStatusEvent(v)
	}

	s.remove(ctx, e)
	s.emails = append(s.emails, newE)
	return newE, true, nil
}
