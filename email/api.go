package email

import (
	"fmt"

	"github.com/husainaloos/notfy/messaging"
	"github.com/husainaloos/notfy/status"
)

type APIInterface interface {
	Queue(Email) (Email, status.Info, error)
}

type API struct {
	p         messaging.Publisher
	s         Storage
	statusAPI status.APIInterface
}

func NewAPI(p messaging.Publisher, s Storage, statusAPI status.APIInterface) *API {
	return &API{p, s, statusAPI}
}

func (api API) Queue(m Email) (Email, status.Info, error) {
	info, err := api.statusAPI.Create(status.Queued)
	if err != nil {
		return Email{}, status.Info{}, fmt.Errorf("failed to create status: %v", err)
	}
	email, err := s.Insert(m, info)
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

func (api API) marshal(m Email, i status.Info) ([]bytes, error) {

}
