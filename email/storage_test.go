package email

import (
	"context"
	"reflect"
	"testing"
)

func TestMemoryStorageInsert(t *testing.T) {
	type emailArg struct {
		id      int
		from    string
		to      string
		subject string
		body    string
	}
	tests := []struct {
		desc    string
		email   emailArg
		expect  emailArg
		wantErr bool
	}{
		{
			desc:    "should insert email successfull",
			email:   emailArg{0, "james@example.com", "john@example.com", "subject", "body"},
			expect:  emailArg{1, "james@example.com", "john@example.com", "subject", "body"},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			s := NewMemoryStorage()
			email, err := New(test.email.id, test.email.from, []string{test.email.to}, nil, nil, test.email.subject, test.email.body)
			if err != nil {
				t.Errorf("%s: failed to create email: %v", test.desc, err)
			}
			expect, err := New(test.expect.id, test.expect.from, []string{test.expect.to}, nil, nil, test.expect.subject, test.expect.body)
			if err != nil {
				t.Errorf("%s: failed to create email: %v", test.desc, err)
			}

			got, err := s.insert(context.Background(), email)
			if test.wantErr && err == nil {
				t.Errorf("%s: expected error, but got no error", test.desc)
			}
			if !test.wantErr && err != nil {
				t.Errorf("%s: got error %v, but expected no error", test.desc, err)
			}
			if !reflect.DeepEqual(got, expect) {
				t.Errorf("%s: got %v, but expected %v", test.desc, got, expect)
			}
		})
	}
}
