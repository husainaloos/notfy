package email

import (
	"net/mail"
	"reflect"
	"testing"
)

func TestNewEmail(t *testing.T) {
	addr1, _ := mail.ParseAddress("email1@gmail.com")
	namedAddr1, _ := mail.ParseAddress("user<email1@gmail.com>")
	addr2, _ := mail.ParseAddress("email2@gmail.com")
	addr, _ := mail.ParseAddress("email@gmail.com")
	unnamedEmail := Email{
		from:    addr,
		to:      []*mail.Address{addr1, addr2},
		cc:      []*mail.Address{addr1, addr2},
		bcc:     []*mail.Address{addr1, addr2},
		subject: "subject",
		body:    "body",
	}
	namedEmail := Email{
		from:    addr,
		to:      []*mail.Address{namedAddr1, addr2},
		cc:      []*mail.Address{namedAddr1, addr2},
		bcc:     []*mail.Address{namedAddr1, addr2},
		subject: "subject",
		body:    "body",
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
				bcc:     []string{"email1@gmail.com", "email2@gmail.com"},
				cc:      []string{"email1@gmail.com", "email2@gmail.com"},
				to:      []string{"email1@gmail.com", "email2@gmail.com"},
				from:    "email@gmail.com",
				body:    "body",
				subject: "subject",
			},
			want:    unnamedEmail,
			wantErr: false,
		},
		{
			name: "should parse named email correctly",
			args: args{
				bcc:     []string{"user <email1@gmail.com>", "email2@gmail.com"},
				cc:      []string{"user <email1@gmail.com>", "email2@gmail.com"},
				to:      []string{"user <email1@gmail.com>", "email2@gmail.com"},
				from:    "email@gmail.com",
				body:    "body",
				subject: "subject",
			},
			want:    namedEmail,
			wantErr: false,
		},
		{
			name: "should fail if from is malformed",
			args: args{
				bcc:     []string{"email1@gmail.com", "email2@gmail.com"},
				cc:      []string{"email1@gmail.com", "email2@gmail.com"},
				to:      []string{"email1@gmail.com", "email2@gmail.com"},
				from:    "emailgmail.com",
				body:    "body",
				subject: "subject",
			},
			want:    Email{},
			wantErr: true,
		},
		{
			name: "should fail if to is malformed",
			args: args{
				bcc:     []string{"email1@gmail.com", "email2@gmail.com"},
				cc:      []string{"email1@gmail.com", "email2@gmail.com"},
				to:      []string{"email1gmail.com", "email2@gmail.com"},
				from:    "email@gmail.com",
				body:    "body",
				subject: "subject",
			},
			want:    Email{},
			wantErr: true,
		},
		{
			name: "should fail if cc is malformed",
			args: args{
				bcc:     []string{"email1@gmail.com", "email2@gmail.com"},
				cc:      []string{"email1gmail.com", "email2@gmail.com"},
				to:      []string{"email1@gmail.com", "email2@gmail.com"},
				from:    "email@gmail.com",
				body:    "body",
				subject: "subject",
			},
			want:    Email{},
			wantErr: true,
		},
		{
			name: "should fail if bcc is malformed",
			args: args{
				bcc:     []string{"email1gmail.com", "email2@gmail.com"},
				cc:      []string{"email1@gmail.com", "email2@gmail.com"},
				to:      []string{"email1@gmail.com", "email2@gmail.com"},
				from:    "email@gmail.com",
				body:    "body",
				subject: "subject",
			},
			want:    Email{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.from, tt.args.to, tt.args.cc, tt.args.bcc, tt.args.subject, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
