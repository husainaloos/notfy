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
	consumers                []messaging.Subscriber
	storage                  Storage
	addr, username, password string
	nclients                 int
	clients                  chan *Client
}

func NewDeamon(consumers []messaging.Subscriber, storage Storage, cfg DeamonConfig) *Deamon {
	clients := make(chan *Client, cfg.SMTPConnectionCount)
	d := &Deamon{
		consumers: consumers,
		storage:   storage,
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
	d.processMessages(ctx, msgC)
}

func (d *Deamon) consume(ctx context.Context) chan []byte {
	msgC := make(chan []byte)
	sf := func(b []byte) {
		msgC <- b
	}
	go func() {
		for _, c := range d.consumers {
			c.Subscribe(sf)
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

func (d *Deamon) processMessages(ctx context.Context, msgC chan []byte) {
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
			logger := logrus.WithField("email_id", email.ID())
			logger.Info("email received")
			emailSent := false
			for i := 0; i < 5 && !emailSent; i++ {
				countLogger := logger.WithField("retry_count", i+1)
				countLogger.Debug("trying to send email")
				if err := c.Send(email); err != nil {
					countLogger.Errorf("failed to send email: %v", err)
					email.AddStatusEvent(MakeStatusEvent(FailedAttemptToSend, time.Now()))
					time.Sleep(5 * time.Second)
					c = d.recycleClient(c)
					continue
				}
				countLogger.Info("email sent")
				emailSent = true

			}
			if emailSent {
				email.AddStatusEvent(MakeStatusEvent(SentSuccessfully, time.Now()))
			} else {
				logger.Error("email is dead")
				email.AddStatusEvent(MakeStatusEvent(Dead, time.Now()))
			}
			d.putClient(c)
			_, ok, err := d.storage.update(ctx, email)
			if err != nil {
				logger.Errorf("failed to update email: %v", err)
			} else if !ok {
				logger.Errorf("email to update does not exist")
			} else {
				logger.Debug("email updated successfully")
			}
		}(msg, c)
	}
	wg.Wait()
}
