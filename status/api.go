package status

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrStatusNotFound is the error used when the status is not found when creating the status
	ErrStatusNotFound = errors.New("status not found")
)

// API is the api for dealing with status
type API struct {
	s Storage
}

// NewAPI creates a new status API
func NewAPI(s Storage) *API { return &API{s} }

// Create status in repository
func (api *API) Create(ctx context.Context, s SendStatus) (Info, error) {
	i := MakeInfo(0, s, time.Now(), time.Now())
	return api.s.insert(ctx, i)
}

// Get status from repository
func (api *API) Get(ctx context.Context, id int) (Info, error) {
	res, err := api.s.get(ctx, id)
	if err != nil {
		switch err {
		case errStorageNotFound:
			return res, ErrStatusNotFound
		default:
			return res, err
		}
	}
	return res, nil

}
