package status

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestInMemoryStorageGet(t *testing.T) {
	s := NewInMemoryStorage()
	ctx := context.Background()
	now := time.Now()
	info1, _ := s.insert(ctx, MakeInfo(0, Sent, now, now))
	info2, _ := s.insert(ctx, MakeInfo(0, Failed, now, now))
	badID := info2.ID() + 1

	expectedSent := MakeInfo(info1.ID(), Sent, now, now)
	expectedFailed := MakeInfo(info2.ID(), Failed, now, now)

	tt := []struct {
		name     string
		id       int
		expected Info
		wantErr  bool
	}{
		{
			name:     "when given id of sent status, should return the sent status",
			id:       info1.ID(),
			expected: expectedSent,
			wantErr:  false,
		},
		{
			name:     "when given id of failed status, should return failed status",
			id:       info2.ID(),
			expected: expectedFailed,
			wantErr:  false,
		},
		{
			name:     "should return NotFoundErr if the id does not exists",
			id:       badID,
			expected: Info{},
			wantErr:  true,
		},
	}

	for _, tst := range tt {
		t.Run(tst.name, func(t *testing.T) {
			ctx := context.Background()
			res, err := s.get(ctx, tst.id)
			if tst.wantErr && err == nil {
				t.Error("expected an error, but found none")
			} else if !tst.wantErr && err != nil {
				t.Errorf("expected no error, but found %v", err)
			} else if tst.expected.Status() != res.Status() {
				t.Errorf("expected %v, but found %v", tst.expected, res)
			}
		})
	}
}

func TestInMemoryStorageUpdate(t *testing.T) {
	s := NewInMemoryStorage()
	ctx := context.Background()
	now := time.Now()
	info, _ := s.insert(ctx, MakeInfo(0, Sent, now, now))

	ue := MakeInfo(info.ID(), Queued, now, now)
	nue := MakeInfo(info.ID()+1, Queued, now, now)
	r := ue

	tt := []struct {
		name         string
		updateEntity Info
		result       Info
		wantErr      bool
	}{
		{
			name:         "when given existing id, the StatusInfo should be updated",
			updateEntity: ue,
			result:       r,
			wantErr:      false,
		},
		{
			name:         "when given non-existing id the StatusInfo should return ErrNotFound",
			updateEntity: nue,
			result:       r,
			wantErr:      true,
		},
	}
	for _, tst := range tt {
		t.Run(tst.name, func(t *testing.T) {
			ctx := context.Background()
			r, err := s.update(ctx, tst.updateEntity)
			if tst.wantErr && err == nil {
				t.Error("expected an error, but found none")
			} else if !tst.wantErr && err != nil {
				t.Errorf("expected no error, but found %v", err)
			} else if !tst.wantErr && !reflect.DeepEqual(tst.result, r) {
				t.Errorf("expected %v, but found %v", tst.result, r)
			}
		})
	}
}
