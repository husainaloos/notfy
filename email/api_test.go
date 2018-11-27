package email

import (
	"reflect"
	"testing"

	"github.com/husainaloos/notfy/status"

	"github.com/husainaloos/notfy/messaging"
)

func Test_QueueEmail(t *testing.T) {
	passCreate := func(s status.SendStatus) (status.Info, error) { return status.MakeInfo(1, s), nil }
	passGet := func(int) (status.Info, error) { return status.Info{}, nil }
	email, _ := New("my@gmail.com", []string{"you@gmail.com"}, []string{}, []string{}, "subject", "body")

	t.Run("should queue an email in the broker", func(t *testing.T) {
		broker := messaging.NewInMemoryPublisher()
		storage := NewInMemoryStorage()
		statusAPI := status.NewMockAPI(passCreate, passGet)
		api := NewAPI(broker, storage, statusAPI)
		e, info, err := api.Queue(email)
		if err != nil {
			t.Errorf("expected no error, but got one: %v", err)
		}
		if info.Status() != status.Queued {
			t.Errorf("info.Status() should be %v, but found %v", status.Queued, info.Status())
		}
		if e.ID() == 0 {
			t.Error("expected e.ID() to be greated that 0, but found 0")
		}
	})
}

func Test_Get(t *testing.T) {
	passCreate := func(status.SendStatus) (status.Info, error) { return status.Info{}, nil }

	t.Run("should return email and status", func(t *testing.T) {
		emailID := 1
		infoID := 2
		expEmail, _ := New("my@gmail.com", []string{"you@gmail.com"}, []string{}, []string{}, "subject", "body")
		expEmail.SetID(1)
		expInfo := status.MakeInfo(infoID, status.Queued)
		passGet := func(int) (status.Info, error) { return expInfo, nil }
		broker := messaging.NewInMemoryPublisher()
		storage := NewInMemoryStorage()
		storage.insert(expEmail, expInfo)
		statusAPI := status.NewMockAPI(passCreate, passGet)
		api := NewAPI(broker, storage, statusAPI)
		e, info, err := api.Get(emailID)
		if err != nil {
			t.Errorf("expected no error, but got one: %v", err)
		}
		if !reflect.DeepEqual(e, expEmail) {
			t.Errorf("expected email %v, but found email %v", expEmail, e)
		}
		if !reflect.DeepEqual(info, expInfo) {
			t.Errorf("expected info %v, but found info %v", expInfo, info)
		}
	})
}
