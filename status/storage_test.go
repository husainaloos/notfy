package status

import (
	"reflect"
	"testing"
)

func Test_InMemoryStorage_Get(t *testing.T) {
	s := NewInMemoryStorage()
	info1, _ := s.insert(MakeInfo(0, Sent))
	info2, _ := s.insert(MakeInfo(0, Failed))
	badID := info2.ID() + 1

	expectedSent := MakeInfo(info1.ID(), Sent)
	expectedFailed := MakeInfo(info2.ID(), Failed)

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
			res, err := s.get(tst.id)
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

func Test_InMemoryStorage_Update(t *testing.T) {
	s := NewInMemoryStorage()
	info, _ := s.insert(MakeInfo(0, Sent))

	ue := MakeInfo(info.ID(), Queued)
	nue := MakeInfo(info.ID()+1, Queued)
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
			r, err := s.update(tst.updateEntity)
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
