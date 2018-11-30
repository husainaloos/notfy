package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/husainaloos/notfy/dto"
	"github.com/husainaloos/notfy/messaging"
	"github.com/sirupsen/logrus"
)

type Deamon struct {
	client   *smtp.Client
	consumer messaging.Consumer
}

func NewDeamon(addr string, username, password string, consumer messaging.Consumer) (*Deamon, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("cannot build host: %v", err)
	}
	auth := smtp.PlainAuth("", username, password, host)
	client, err := smtp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial given addr: %v", err)
	}
	config := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}
	if err := client.StartTLS(config); err != nil {
		return nil, fmt.Errorf("cannot start TLS connection with smtp server: %v", err)
	}
	if err := client.Auth(auth); err != nil {
		return nil, fmt.Errorf("failed to create auth: %v", err)
	}
	return &Deamon{
		client:   client,
		consumer: consumer,
	}, nil
}

func (d *Deamon) Start() error {
	logrus.Debug("deamon instantiating...")
	for {
		b, err := d.consumer.Consume()
		if err != nil {
			return fmt.Errorf("failed to consumer: %v", err)
		}
		logrus.WithField("msg_size", len(b)).Debug("consumed message")
		var pe dto.PublishedEmail
		if err := proto.Unmarshal(b, &pe); err != nil {
			return fmt.Errorf("failed to unmarshal message: %v", err)
		}

		st := time.Now()
		if err := d.sendEmail(pe); err != nil {
			return fmt.Errorf("failed to send email: %v", err)
		}
		d := time.Since(st)
		logrus.WithFields(logrus.Fields{
			"elapsed":  float64(d) / 1000000.0,
			"msg_size": len(b),
		}).Info("email sent")
	}
}

func (d *Deamon) sendEmail(pe dto.PublishedEmail) error {
	logrus.WithField("from", pe.From).Debug("processing from")
	if err := d.client.Mail(pe.From); err != nil {
		return err
	}

	if pe.To == nil {
		pe.To = make([]string, 0)
	}
	logrus.WithField("to_size", len(pe.To)).Debug("processing to")
	for _, r := range pe.To {
		if err := d.client.Rcpt(r); err != nil {
			return err
		}
	}

	if pe.Cc == nil {
		pe.Cc = make([]string, 0)
	}
	logrus.WithField("cc_size", len(pe.Cc)).Debug("processing cc")
	for _, r := range pe.Cc {
		if err := d.client.Rcpt(r); err != nil {
			return err
		}
	}

	if pe.Bcc == nil {
		pe.Bcc = make([]string, 0)
	}
	logrus.WithField("bcc_size", len(pe.Bcc)).Debug("processing bcc")
	for _, r := range pe.Bcc {
		if err := d.client.Rcpt(r); err != nil {
			return err
		}
	}

	logrus.Debug("building message body")
	wc, err := d.client.Data()
	if err != nil {
		return err
	}

	sb := strings.Builder{}
	sb.WriteString("From: ")
	sb.WriteString(pe.From)
	sb.WriteString("\r\n")
	sb.WriteString("To:")
	sb.WriteString(strings.Join(pe.To, ","))
	sb.WriteString("\r\n")
	sb.WriteString("Cc:")
	sb.WriteString(strings.Join(pe.Cc, ","))
	sb.WriteString("\r\n")
	sb.WriteString("Bcc:")
	sb.WriteString(strings.Join(pe.Bcc, ","))
	sb.WriteString("\r\n")
	sb.WriteString("Subject:")
	sb.WriteString(pe.Subject)
	sb.WriteString("\r\n\r\n")
	sb.WriteString(pe.Body)
	wc.Write([]byte(sb.String()))
	return wc.Close()
}
