package status

import "time"

type APIInterface interface {
	Create(SendStatus) (Info, error)
	Get(id int) (Info, error)
}

// API is the api for dealing with status
type API struct {
	s Storage
}

// NewAPI creates a new status API
func NewAPI(s Storage) *API { return &API{s} }

// Create status in repository
func (api *API) Create(s SendStatus) (Info, error) {
	i := MakeInfo(0, s, time.Now(), time.Now())
	return api.s.Insert(i)
}

// Get status from repository
func (api *API) Get(id int) (Info, error) {
	return api.s.Get(id)
}
