package email

import (
	"encoding/json"
	"errors"
	"net/mail"
	"time"
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
func (m Email) ID() int       { return m.id }
func (m *Email) SetID(id int) { m.id = id }

// From gets the from
func (m Email) From() mail.Address {
	return *m.from
}
func (m Email) StringFrom() string {
	return m.from.String()
}

// To gets the to
func (m Email) To() []mail.Address {
	arr := make([]mail.Address, len(m.to))
	for i, v := range m.to {
		arr[i] = *v
	}
	return arr
}

func (m Email) StringTo() []string {
	arr := []string{}
	for _, v := range m.to {
		arr = append(arr, v.String())
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

func (m Email) StringCC() []string {
	arr := []string{}
	for _, v := range m.cc {
		arr = append(arr, v.String())
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

func (m Email) StringBCC() []string {
	arr := []string{}
	for _, v := range m.bcc {
		arr = append(arr, v.String())
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

// StatusHistory gets the status history of the email
func (m Email) StatusHistory() StatusHistory {
	sh := make(StatusHistory, 0)
	for _, v := range m.statusHistory {
		sh = append(sh, v)
	}
	return sh
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

func (e Email) testString() string {
	type se struct {
		Status int       `json:"status"`
		At     time.Time `json:"at"`
	}
	type testEmail struct {
		ID           int      `json:"id"`
		From         string   `json:"from"`
		To           []string `json:"to"`
		CC           []string `json:"cc"`
		BCC          []string `json:"bcc"`
		Subject      string   `json:"subject"`
		Body         string   `json:"body"`
		StatusEvents []se     `json:"status_events"`
	}

	te := testEmail{
		ID:           e.ID(),
		From:         e.StringFrom(),
		To:           e.StringTo(),
		CC:           e.StringCC(),
		BCC:          e.StringBCC(),
		Subject:      e.Subject(),
		Body:         e.Body(),
		StatusEvents: []se{},
	}

	for _, v := range e.StatusHistory() {
		te.StatusEvents = append(te.StatusEvents, se{int(v.Status()), v.At()})
	}

	j, _ := json.MarshalIndent(te, "", "  ")
	return string(j)
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
