package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/sirupsen/logrus"
)

type Client struct {
	smtpc *smtp.Client
}

func NewClient(addr string, username, password string) (*Client, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("cannot build host: %v", err)
	}
	auth := smtp.PlainAuth("", username, password, host)
	smtpc, err := smtp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial given addr: %v", err)
	}
	config := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}
	if err := smtpc.StartTLS(config); err != nil {
		return nil, fmt.Errorf("cannot start TLS connection with smtp server: %v", err)
	}
	if err := smtpc.Auth(auth); err != nil {
		return nil, fmt.Errorf("failed to create auth: %v", err)
	}
	return &Client{
		smtpc: smtpc,
	}, nil

}

func (c *Client) Send(e Email) error {
	maddr := e.From()
	from := maddr.String()
	logrus.WithField("from", from).Debug("processing from")
	if err := c.smtpc.Mail(from); err != nil {
		return err
	}
	to := make([]string, 0)
	logrus.WithField("to_size", len(e.To())).Debug("processing to")
	for _, r := range e.To() {
		maddr := r.String()
		to = append(to, maddr)
		if err := c.smtpc.Rcpt(maddr); err != nil {
			return err
		}
	}
	cc := make([]string, 0)
	logrus.WithField("cc_size", len(e.CC())).Debug("processing cc")
	for _, r := range e.CC() {
		maddr := r.String()
		cc = append(cc, maddr)
		if err := c.smtpc.Rcpt(maddr); err != nil {
			return err
		}
	}
	bcc := make([]string, 0)
	logrus.WithField("bcc_size", len(e.BCC())).Debug("processing bcc")
	for _, r := range e.BCC() {
		maddr := r.String()
		bcc = append(bcc, maddr)
		if err := c.smtpc.Rcpt(maddr); err != nil {
			return err
		}
	}
	logrus.Debug("building message body")
	wc, err := c.smtpc.Data()
	if err != nil {
		return err
	}
	sb := strings.Builder{}
	sb.WriteString("From: ")
	sb.WriteString(from)
	sb.WriteString("\r\n")
	sb.WriteString("To:")
	sb.WriteString(strings.Join(to, ","))
	sb.WriteString("\r\n")
	sb.WriteString("Cc:")
	sb.WriteString(strings.Join(cc, ","))
	sb.WriteString("\r\n")
	sb.WriteString("Bcc:")
	sb.WriteString(strings.Join(bcc, ","))
	sb.WriteString("\r\n")
	sb.WriteString("Subject:")
	sb.WriteString(e.Subject())
	sb.WriteString("\r\n\r\n")
	sb.WriteString(e.Body())
	wc.Write([]byte(sb.String()))
	return wc.Close()
}

func (c *Client) Close() error {
	return c.smtpc.Close()
}
