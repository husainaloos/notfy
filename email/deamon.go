package email

import (
	"context"
	"sync"
	"time"

	"github.com/husainaloos/notfy/messaging"
	"github.com/sirupsen/logrus"
)

type DeamonConfig struct {
	SMTPAddr            string
	SMTPUsername        string
	SMTPPassword        string
	SMTPConnectionCount int
}

type Deamon struct {
	consumers                []messaging.Consumer
	addr, username, password string
	nclients                 int
	clients                  chan *Client
}

func NewDeamon(consumers []messaging.Consumer, cfg DeamonConfig) *Deamon {
	clients := make(chan *Client, cfg.SMTPConnectionCount)
	d := &Deamon{
		consumers: consumers,
		addr:      cfg.SMTPAddr,
		username:  cfg.SMTPUsername,
		password:  cfg.SMTPPassword,
		nclients:  cfg.SMTPConnectionCount,
		clients:   clients,
	}

	// generate the smtp clients
	go func(d *Deamon) {
		created := 0
		for created < d.nclients {
			c, err := NewClient(d.addr, d.username, d.password)
			if err != nil {
				logrus.Errorf("cannot create client: %v", err)
				logrus.Info("retry creating client")
				continue
			}
			logrus.Debugf("client %d created", created)
			d.clients <- c
			created++
		}
		logrus.Infof("created %d clients", created)
	}(d)
	return d
}

func (d *Deamon) Start(ctx context.Context) {
	logrus.Debug("deamon starting")
	msgC := d.consume(ctx)
	d.processMessages(msgC)
}

func (d *Deamon) consume(ctx context.Context) chan []byte {
	msgC := make(chan []byte)
	go func() {
		for _, c := range d.consumers {
			go func(c messaging.Consumer, msgC chan []byte) {
				logrus.Debug("starting routine for consumer")
				for {
					select {
					case <-ctx.Done():
						close(msgC)
						logrus.Info("context is cancelled: channel closed")
						return
					default:
						b, err := c.Consume()
						if err != nil {
							logrus.Errorf("failed to receive message from consumer: %v", err)
							continue
						}
						logrus.WithField("msg_size", len(b)).Info("message received from consumer")
						msgC <- b
					}
				}
			}(c, msgC)
		}
	}()
	return msgC
}

// get a client from the client pool
func (d *Deamon) getClient() *Client {
	return <-d.clients
}

// put the client back to the pool
func (d *Deamon) putClient(c *Client) {
	d.clients <- c
}

// rebuild the client
func (d *Deamon) recycleClient(c *Client) *Client {
	go func() {
		for {
			newC, err := NewClient(d.addr, d.username, d.password)
			if err != nil {
				logrus.Errorf("cannot create client: %v", err)
				continue
			}
			d.clients <- newC
		}
	}()
	c.Close()
	return <-d.clients
}

func (d *Deamon) processMessages(msgC chan []byte) {
	logrus.Debug("email sending routine started")
	var wg sync.WaitGroup
	for msg := range msgC {
		logrus.WithField("msg_size", len(msg)).Debug("message about to be send")
		wg.Add(1)
		c := d.getClient()
		go func(msg []byte, c *Client) {
			defer wg.Done()
			email, err := Unmarshal(msg)
			if err != nil {
				logrus.Errorf("cannot parse email: %v", err)
				return
			}
			logrus.WithField("email_id", email.ID()).Info("email received")
			emailSent := false
			for i := 0; i < 5 && !emailSent; i++ {
				logrus.WithFields(logrus.Fields{
					"email_id":    email.ID(),
					"retry_count": i + 1,
				}).Debug("trying to send email")
				if err := c.Send(email); err != nil {
					logrus.WithFields(logrus.Fields{
						"email_id":    email.ID(),
						"retry_count": i + 1,
					}).Errorf("failed to send email: %v", err)
					email.AddStatusEvent(MakeStatusEvent(FailedAttemptToSend, time.Now()))
					c = d.recycleClient(c)
					continue
				}
				logrus.WithField("email_id", email.ID()).Info("email sent")
				emailSent = true
			}
			if emailSent {
				email.AddStatusEvent(MakeStatusEvent(SentSuccessfully, time.Now()))
			} else {
				logrus.WithField("email_id", email.ID()).Error("email is dead")
				email.AddStatusEvent(MakeStatusEvent(Dead, time.Now()))
			}
			d.putClient(c)
		}(msg, c)
	}
	wg.Wait()
}
