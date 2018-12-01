package status

import (
	"context"
	"testing"
	"time"
)

func TestAPICreate(t *testing.T) {
	tt := []struct {
		name     string
		status   SendStatus
		expected Info
		wantErr  bool
	}{
		{
			name:     "should insert status to storage",
			status:   Queued,
			expected: MakeInfo(1, Queued, time.Now(), time.Now()),
			wantErr:  false,
		},
	}
	for _, tst := range tt {
		t.Run(tst.name, func(t *testing.T) {
			storage := NewInMemoryStorage()
			api := NewAPI(storage)
			got, err := api.Create(context.Background(), tst.status)
			if tst.wantErr && err == nil {
				t.Errorf("Create(): got no error, but wanted an error")
			}
			if !tst.wantErr && err != nil {
				t.Errorf("Create(): got error %v, but wanted no error", err)
			}
			if got.Status() != tst.expected.Status() {
				t.Errorf("Create(): got status %v, but expected %v", got.Status(), tst.expected.Status())
			}
		})
	}
}
