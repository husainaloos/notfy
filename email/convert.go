package email

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/husainaloos/notfy/dto"
	"github.com/sirupsen/logrus"
)

func Marshal(e Email) ([]byte, error) {
	p := &dto.QueuedEmail{
		Id:      uint64(e.ID()),
		Subject: e.Subject(),
		Body:    e.Body(),
	}
	from := e.From()
	to := []string{}
	cc := []string{}
	bcc := []string{}
	se := []*dto.StatusEvent{}
	for _, v := range e.To() {
		to = append(to, v.String())
	}
	for _, v := range e.CC() {
		cc = append(cc, v.String())
	}
	for _, v := range e.BCC() {
		bcc = append(bcc, v.String())
	}
	for _, v := range e.StatusHistory() {
		s := &dto.StatusEvent{
			Status: uint32(v.Status()),
			At:     uint64(v.At().UnixNano()),
		}
		se = append(se, s)
	}

	p.To = to
	p.Cc = cc
	p.Bcc = bcc
	p.From = from.String()
	p.Status = se

	logrus.Debug("email to queue", p)
	return proto.Marshal(p)
}

func Unmarshal(b []byte) (Email, error) {
	p := &dto.QueuedEmail{}
	err := proto.Unmarshal(b, p)
	if err != nil {
		return Email{}, err
	}
	e, err := New(int(p.Id), p.From, p.To, p.Cc, p.Bcc, p.Subject, p.Body)
	if err != nil {
		return Email{}, err
	}
	for _, v := range p.Status {
		s := Status(v.Status)
		t := time.Unix(0, int64(v.At))
		se := MakeStatusEvent(s, t)
		e.AddStatusEvent(se)
	}
	return e, nil
}
