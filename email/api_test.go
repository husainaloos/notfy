package email

import (
	"testing"

	"github.com/husainaloos/notfy/status"

	"github.com/husainaloos/notfy/messaging"
)

func Test_QueueEmail(t *testing.T) {
	t.Run("should queue an email in the broker", func(t *testing.T) {
		email, _ := New("my@gmail.com", []string{"you@gmail.com"}, []string{}, []string{}, "subject", "body")
		broker := messaging.NewInMemoryPublisher()
		storage := NewInMemoryStorage()
		statusAPI := status.NewAPI(status.NewInMemoryStorage())
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
