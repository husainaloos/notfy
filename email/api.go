package email

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/husainaloos/notfy/dto"
	"github.com/husainaloos/notfy/messaging"
	"github.com/husainaloos/notfy/status"
)

// APIInterface is an interface for the API
type APIInterface interface {
	Queue(Email) (Email, status.Info, error)
	Get(id int) (Email, status.Info, error)
}

// API is the api for dealing with emails
type API struct {
	p         messaging.Publisher
	s         Storage
	statusAPI status.APIInterface
}

// NewAPI creates a new API
func NewAPI(p messaging.Publisher, s Storage, statusAPI status.APIInterface) *API {
	return &API{p, s, statusAPI}
}

// Queue an email to be sent
func (api API) Queue(m Email) (Email, status.Info, error) {
	info, err := api.statusAPI.Create(status.Queued)
	if err != nil {
		return Email{}, status.Info{}, fmt.Errorf("failed to create status: %v", err)
	}
	email, err := api.s.insert(m, info)
	if err != nil {
		return Email{}, status.Info{}, fmt.Errorf("failed to insert an email to the storage: %v", err)
	}
	b, err := api.marshal(email, info)
	if err != nil {
		return Email{}, status.Info{}, fmt.Errorf("failed to marshal email to binary: %v", err)
	}
	if err := api.p.Publish(b); err != nil {
		return Email{}, status.Info{}, fmt.Errorf("failed to publish the email to the publisher: %v", err)
	}
	return email, info, nil
}

// Get the email info
func (api API) Get(id int) (Email, status.Info, error) {
	se, err := api.s.get(id)
	if err != nil {
		return Email{}, status.Info{}, fmt.Errorf("cannot get email from storage: %v", err)
	}
	info, err := api.statusAPI.Get(se.statusID)
	if err != nil {
		return Email{}, status.Info{}, fmt.Errorf("cannot get status of email: %v", err)
	}
	return se.Email, info, nil
}

func (api API) marshal(m Email, i status.Info) ([]byte, error) {
	pe := &dto.PublishedEmail{
		From:        m.from.String(),
		To:          make([]string, 0),
		Cc:          make([]string, 0),
		Bcc:         make([]string, 0),
		Subject:     m.Subject(),
		Body:        m.Body(),
		Id:          int64(m.ID()),
		StatusId:    int64(i.ID()),
		PublishedAt: time.Now().Unix(),
	}
	for _, v := range m.To() {
		pe.To = append(pe.To, v.String())
	}
	for _, v := range m.CC() {
		pe.Cc = append(pe.Cc, v.String())
	}
	for _, v := range m.BCC() {
		pe.Bcc = append(pe.Bcc, v.String())
	}
	b, err := proto.Marshal(pe)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// MockAPI is a mock of APIInterface
type MockAPI struct {
	queue func(Email) (Email, status.Info, error)
	get   func(int) (Email, status.Info, error)
}

// NewMockAPI creates a new mock
func NewMockAPI(
	queue func(Email) (Email, status.Info, error),
	get func(int) (Email, status.Info, error)) *MockAPI {
	return &MockAPI{queue, get}
}

// Queue mocked
func (api *MockAPI) Queue(e Email) (Email, status.Info, error) {
	return api.queue(e)
}

// Get mocked
func (api *MockAPI) Get(id int) (Email, status.Info, error) {
	return api.get(id)
}
