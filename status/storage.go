package status

import (
	"errors"
	"sync"
)

var (
	errStorageNotFound = errors.New("record not found in storage")
)

// Storage is an interface for storing Info
type Storage interface {
	insert(Info) (Info, error)
	update(Info) (Info, error)
	get(id int) (Info, error)
}

// InMemoryStorage stores Info in memory
type InMemoryStorage struct {
	infos []Info
	mu    sync.Mutex
	id    int
}

// NewInMemoryStorage creates new InMemoryStorage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		infos: make([]Info, 0),
		id:    0,
	}
}

func (s *InMemoryStorage) insert(r Info) (Info, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.id++
	r.id = s.id
	s.infos = append(s.infos, r)

	return r, nil
}

func (s *InMemoryStorage) update(r Info) (Info, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, si := range s.infos {
		if si.ID() == r.ID() {
			si.createdAt = r.createdAt
			si.lastUpdateAt = r.lastUpdateAt
			si.status = r.status
			return si, nil
		}
	}
	return Info{}, errStorageNotFound
}

func (s *InMemoryStorage) get(id int) (Info, error) {
	for _, r := range s.infos {
		if r.ID() == id {
			return r, nil
		}
	}
	return Info{}, errStorageNotFound
}
