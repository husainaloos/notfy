package status

import (
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
func (api *API) Create(s SendStatus) (Info, error) {
	i := MakeInfo(0, s)
	i.SetCreatedAt(time.Now())
	i.SetLastUpdatedAt(time.Now())
	return api.s.insert(i)
}

// Get status from repository
func (api *API) Get(id int) (Info, error) {
	res, err := api.s.get(id)
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
