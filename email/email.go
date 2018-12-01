package email

import (
	"errors"
	"net/mail"
)

// Email is the email struct
type Email struct {
	id            int
	from          *mail.Address
	to            []*mail.Address
	cc            []*mail.Address
	bcc           []*mail.Address
	subject       string
	body          string
	statusHistory StatusHistory
}

// ID gets the id of the email
func (m Email) ID() int { return m.id }

// From gets the from
func (m Email) From() mail.Address {
	return *m.from
}

// To gets the to
func (m Email) To() []mail.Address {
	arr := make([]mail.Address, len(m.to))
	for i, v := range m.to {
		arr[i] = *v
	}
	return arr
}

// CC gets the cc
func (m Email) CC() []mail.Address {
	arr := make([]mail.Address, len(m.cc))
	for i, v := range m.cc {
		arr[i] = *v
	}
	return arr
}

// BCC gets the bcc
func (m Email) BCC() []mail.Address {
	arr := make([]mail.Address, len(m.bcc))
	for i, v := range m.bcc {
		arr[i] = *v
	}
	return arr
}

// Subject gets the subject
func (m Email) Subject() string {
	return m.subject
}

// Body gets the body
func (m Email) Body() string {
	return m.body
}

// AddStatusEvent adds an event to the email
func (m *Email) AddStatusEvent(se StatusEvent) {
	m.statusHistory = append(m.statusHistory, se)
}

// New creates an email
func New(id int, from string, to, cc, bcc []string, subject, body string) (Email, error) {
	if from == "" {
		return Email{}, errors.New("from cannot be empty")
	}
	if cc == nil {
		cc = make([]string, 0)
	}
	if bcc == nil {
		bcc = make([]string, 0)
	}
	if to == nil {
		to = make([]string, 0)
	}
	if len(to) == 0 && len(cc) == 0 && len(bcc) == 0 {
		return Email{}, errors.New("email should have at least one recepient")
	}
	f, err := mail.ParseAddress(from)
	if err != nil {
		return Email{}, err
	}
	tos, err := parseAddList(to)
	if err != nil {
		return Email{}, err
	}
	ccs, err := parseAddList(cc)
	if err != nil {
		return Email{}, err
	}
	bccs, err := parseAddList(bcc)
	if err != nil {
		return Email{}, err
	}
	return Email{id, f, tos, ccs, bccs, subject, body, make(StatusHistory, 0)}, nil
}

func parseAddList(addrs []string) ([]*mail.Address, error) {
	arr := make([]*mail.Address, 0)
	for _, v := range addrs {
		a, err := mail.ParseAddress(v)
		if err != nil {
			return nil, err
		}
		arr = append(arr, a)
	}
	return arr, nil
}
