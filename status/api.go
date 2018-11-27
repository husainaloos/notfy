package status

import "time"

// APIInterface is the interface for the API
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
	i := MakeInfo(0, s)
	i.SetCreatedAt(time.Now())
	i.SetLastUpdatedAt(time.Now())
	return api.s.insert(i)
}

// Get status from repository
func (api *API) Get(id int) (Info, error) {
	return api.s.get(id)
}

// MockAPI is a mock of the API
type MockAPI struct {
	create func(SendStatus) (Info, error)
	get    func(int) (Info, error)
}

// NewMockAPI creates a new mock
func NewMockAPI(create func(SendStatus) (Info, error), get func(int) (Info, error)) *MockAPI {
	return &MockAPI{create, get}
}

// Create mocked
func (api *MockAPI) Create(s SendStatus) (Info, error) {
	return api.create(s)
}

// Get mocked
func (api *MockAPI) Get(id int) (Info, error) {
	return api.get(id)
}
