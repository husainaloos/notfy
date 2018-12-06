package email

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/husainaloos/notfy/dto"
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
	b, err := api.marshal(e)
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

func (api *API) marshal(e Email) ([]byte, error) {
	p := &dto.PublishedEmail{
		Id:      int64(e.ID()),
		Subject: e.Subject(),
		Body:    e.Body(),
	}
	from := e.From()
	to := []string{}
	cc := []string{}
	bcc := []string{}
	for _, v := range e.To() {
		to = append(to, v.String())
	}
	for _, v := range e.CC() {
		cc = append(cc, v.String())
	}
	for _, v := range e.BCC() {
		bcc = append(bcc, v.String())
	}
	p.To = to
	p.Cc = cc
	p.Bcc = bcc
	p.From = from.String()

	return proto.Marshal(p)
}
