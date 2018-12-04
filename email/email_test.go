package email

import (
	"fmt"
	"net/mail"
	"reflect"
	"testing"
)

func TestNewEmail(t *testing.T) {
	jonathan, _ := mail.ParseAddress("jonathan@example.com")
	randy, _ := mail.ParseAddress("randy@example.com")
	sam, _ := mail.ParseAddress("sam@example.com")
	jim, _ := mail.ParseAddress("jim@example.com")
	unnamedEmail := Email{
		from:          sam,
		to:            []*mail.Address{jonathan, randy},
		cc:            []*mail.Address{jonathan, randy},
		bcc:           []*mail.Address{jonathan, randy},
		subject:       "subject",
		body:          "body",
		statusHistory: make(StatusHistory, 0),
	}
	namedEmail := Email{
		from:          sam,
		to:            []*mail.Address{jim, randy},
		cc:            []*mail.Address{jim, randy},
		bcc:           []*mail.Address{jim, randy},
		subject:       "subject",
		body:          "body",
		statusHistory: make(StatusHistory, 0),
	}
	type args struct {
		from    string
		to      []string
		cc      []string
		bcc     []string
		subject string
		body    string
	}
	tests := []struct {
		name    string
		args    args
		want    Email
		wantErr bool
	}{
		{
			name: "should parse email correctly",
			args: args{
				bcc:     []string{"jonathan@example.com", "randy@example.com"},
				cc:      []string{"jonathan@example.com", "randy@example.com"},
				to:      []string{"jonathan@example.com", "randy@example.com"},
				from:    "sam@example.com",
				body:    "body",
				subject: "subject",
			},
			want:    unnamedEmail,
			wantErr: false,
		},
		{
			name: "should parse named email correctly",
			args: args{
				bcc:     []string{"Jim <jim@example.com>", "randy@example.com"},
				cc:      []string{"Jim <jim@example.com>", "randy@example.com"},
				to:      []string{"Jim <jim@example.com>", "randy@example.com"},
				from:    "sam@example.com",
				body:    "body",
				subject: "subject",
			},
			want:    namedEmail,
			wantErr: false,
		},
		{
			name: "should fail if from is malformed",
			args: args{
				bcc:     []string{"jonathan@example.com", "randy@example.com"},
				cc:      []string{"jonathan@example.com", "randy@example.com"},
				to:      []string{"jonathan@example.com", "randy@example.com"},
				from:    "emailexample.com",
				body:    "body",
				subject: "subject",
			},
			want:    Email{},
			wantErr: true,
		},
		{
			name: "should fail if to is malformed",
			args: args{
				bcc:     []string{"jonathan@example.com", "randy@example.com"},
				cc:      []string{"jonathan@example.com", "randy@example.com"},
				to:      []string{"email1example.com", "randy@example.com"},
				from:    "sam@example.com",
				body:    "body",
				subject: "subject",
			},
			want:    Email{},
			wantErr: true,
		},
		{
			name: "should fail if cc is malformed",
			args: args{
				bcc:     []string{"jonathan@example.com", "randy@example.com"},
				cc:      []string{"email1example.com", "randy@example.com"},
				to:      []string{"jonathan@example.com", "randy@example.com"},
				from:    "sam@example.com",
				body:    "body",
				subject: "subject",
			},
			want:    Email{},
			wantErr: true,
		},
		{
			name: "should fail if bcc is malformed",
			args: args{
				bcc:     []string{"email1example.com", "randy@example.com"},
				cc:      []string{"jonathan@example.com", "randy@example.com"},
				to:      []string{"jonathan@example.com", "randy@example.com"},
				from:    "sam@example.com",
				body:    "body",
				subject: "subject",
			},
			want:    Email{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(0, tt.args.from, tt.args.to, tt.args.cc, tt.args.bcc, tt.args.subject, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				fmt.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
