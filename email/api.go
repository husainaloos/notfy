package email

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/husainaloos/notfy/messaging"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

type API struct {
	publisher messaging.Publisher
	storage   Storage
}

func NewAPI(p messaging.Publisher, s Storage) *API {
	return &API{p, s}
}

func (api *API) Queue(ctx context.Context, e Email) (Email, error) {
	b, err := Marshal(e)
	if err != nil {
		return Email{}, fmt.Errorf("failed to marshal email to protobuffer: %v", err)
	}
	if err := api.publisher.Publish(b); err != nil {
		return Email{}, fmt.Errorf("failed to publish email: %v", err)
	}
	e.AddStatusEvent(MakeStatusEvent(Queued, time.Now()))
	email, err := api.storage.insert(ctx, e)
	if err != nil {
		return Email{}, err
	}
	return email, nil
}

func (api *API) Get(ctx context.Context, id int) (Email, error) {
	e, ok, err := api.storage.get(ctx, id)
	if err != nil {
		return Email{}, fmt.Errorf("failed to get from db: %v", err)
	}
	if !ok {
		return Email{}, ErrItemNotFound
	}
	return e, nil
}

func (api *API) MarkDone(id int) error {
	return nil
}
