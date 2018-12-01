package email

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/husainaloos/notfy/status"

	"github.com/husainaloos/notfy/messaging"
)

type mockStatusAPI struct {
	create func(status.SendStatus) (status.Info, error)
	get    func(int) (status.Info, error)
}

func newMockStatusAPI(
	create func(status.SendStatus) (status.Info, error),
	get func(int) (status.Info, error)) *mockStatusAPI {
	return &mockStatusAPI{create, get}
}

func (api *mockStatusAPI) Create(ctx context.Context, s status.SendStatus) (status.Info, error) {
	return api.create(s)
}
func (api *mockStatusAPI) Get(ctx context.Context, id int) (status.Info, error) { return api.get(id) }

func TestQueueEmail(t *testing.T) {
	passCreate := func(s status.SendStatus) (status.Info, error) {
		return status.MakeInfo(1, s, time.Now(), time.Now()), nil
	}
	failCreate := func(s status.SendStatus) (status.Info, error) {
		return status.MakeInfo(1, s, time.Now(), time.Now()), errors.New("create status failed")
	}
	passGet := func(int) (status.Info, error) { return status.Info{}, nil }
	email, _ := New(0, "my@gmail.com", []string{"you@gmail.com"}, []string{}, []string{}, "subject", "body")

	tt := []struct {
		name          string
		createf       func(status.SendStatus) (status.Info, error)
		getf          func(int) (status.Info, error)
		expectedInfo  status.Info
		expectedEmail Email
		wantErr       bool
	}{
		{
			name:          "should queue an email in the broker",
			createf:       passCreate,
			getf:          passGet,
			expectedInfo:  status.MakeInfo(1, status.Queued, time.Now(), time.Now()),
			expectedEmail: email,
			wantErr:       false,
		},
		{
			name:          "should return err if creating status fails",
			createf:       failCreate,
			getf:          passGet,
			expectedInfo:  status.Info{},
			expectedEmail: Email{},
			wantErr:       true,
		},
	}
	for _, tst := range tt {
		t.Run(tst.name, func(t *testing.T) {
			broker := messaging.NewInMemoryBroker()
			storage := NewInMemoryStorage()
			statusAPI := newMockStatusAPI(tst.createf, tst.getf)
			api := NewAPI(broker, storage, statusAPI)
			e, info, err := api.Queue(context.Background(), email)
			if tst.wantErr {
				if err == nil {
					t.Errorf("Queue(): got no error, but wanted an error")
				}
				return
			}
			if err != nil {
				t.Errorf("Queue(): got err %v, expected none", err)
			}
			if info.ID() != tst.expectedInfo.ID() || info.Status() != tst.expectedInfo.Status() {
				t.Errorf("Queue(): got info %v, but expected %v", info, tst.expectedInfo)
			}
			if !reflect.DeepEqual(e, tst.expectedEmail) {
				t.Errorf("Queue(): got email %v, but expected %v", e, tst.expectedEmail)
			}
		})
	}
}

func TestGet(t *testing.T) {
	passCreate := func(status.SendStatus) (status.Info, error) { return status.Info{}, nil }

	// the setup of the test
	// ensure that email (with ID=1) is associated with status (with ID=2)
	email, _ := New(0, "my@gmail.com", []string{"you@gmail.com"}, []string{}, []string{}, "subject", "body")
	info := status.MakeInfo(2, status.Queued, time.Now(), time.Now())
	storage := NewInMemoryStorage()
	storage.insert(email, info)

	tt := []struct {
		name     string
		createf  func(status.SendStatus) (status.Info, error)
		getf     func(int) (status.Info, error)
		emailID  int
		expInfo  status.Info
		expEmail Email
		wantErr  bool
	}{
		{
			name:    "should return email and status",
			createf: passCreate,
			getf: func(id int) (status.Info, error) {
				if id != 2 {
					return status.Info{}, errors.New("bad id")
				}
				return status.MakeInfo(2, status.Queued, time.Now(), time.Now()), nil
			},
			emailID:  1,
			expInfo:  status.MakeInfo(2, status.Queued, time.Now(), time.Now()),
			expEmail: email,
			wantErr:  false,
		},
	}

	for _, tst := range tt {
		t.Run(tst.name, func(t *testing.T) {
			broker := messaging.NewInMemoryBroker()
			statusAPI := newMockStatusAPI(tst.createf, tst.getf)
			api := NewAPI(broker, storage, statusAPI)
			e, i, err := api.Get(context.Background(), tst.emailID)
			if tst.wantErr {
				if err == nil {
					t.Error("Get(): got no error, but expected on")
				}
				return
			}
			if err != nil {
				t.Errorf("Get(): got %v, but expected no error", err)
			}
			if !reflect.DeepEqual(e, tst.expEmail) {
				t.Errorf("Get(): got email %v, but expected %v", e, tst.expEmail)
			}
			if i.ID() != tst.expInfo.ID() || i.Status() != tst.expInfo.Status() {
				t.Errorf("Get(): got info %v, but expected %v", i, tst.expInfo)
			}
		})
	}
}
