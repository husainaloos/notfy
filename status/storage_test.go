package status

import (
	"reflect"
	"testing"
	"time"
)

func Test_InMemoryStorage_Get(t *testing.T) {
	s := NewInMemoryStorage()
	tn := time.Now()
	info1, _ := s.Insert(MakeInfo(0, Sent, tn, tn))
	info2, _ := s.Insert(MakeInfo(0, Failed, tn, tn))
	badID := info2.ID() + 1

	expectedSent := MakeInfo(info1.ID(), Sent, tn, tn)
	expectedFailed := MakeInfo(info2.ID(), Failed, tn, tn)

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
			res, err := s.Get(tst.id)
			if tst.wantErr && err == nil {
				t.Error("expected an error, but found none")
			} else if !tst.wantErr && err != nil {
				t.Errorf("expected no error, but found %v", err)
			} else if !reflect.DeepEqual(tst.expected, res) {
				t.Errorf("expected %v, but found %v", tst.expected, res)
			}
		})
	}
}

func Test_InMemoryStorage_Update(t *testing.T) {
	s := NewInMemoryStorage()
	now := time.Now()
	future := time.Now().Add(1 * time.Second)
	info, _ := s.Insert(MakeInfo(0, Sent, now, now))

	ue := MakeInfo(info.ID(), Queued, future, future)
	nue := MakeInfo(info.ID()+1, Queued, future, future)
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
			r, err := s.Update(tst.updateEntity)
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
