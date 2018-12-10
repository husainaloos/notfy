package email

import (
	"context"
	"time"

	"github.com/husainaloos/notfy/messaging"
	"github.com/sirupsen/logrus"
)

type Deamon struct {
	consumers                []messaging.Consumer
	addr, username, password string
	clients                  chan *Client
}

func NewDeamon(consumers []messaging.Consumer, addr, username, password string) *Deamon {
	clients := make(chan *Client, 50)
	return &Deamon{consumers, addr, username, password, clients}
}

func (d *Deamon) Start(ctx context.Context) {
	logrus.Debug("deamon starting")
	msgC := make(chan []byte)
	go d.startConsuming(ctx, msgC)
	go d.startSending(ctx, msgC)
	<-ctx.Done()
}

func (d *Deamon) startConsuming(ctx context.Context, msgC chan []byte) {
	for _, c := range d.consumers {
		go func(c messaging.Consumer) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					b, err := c.Consume()
					logrus.WithField("msg_size", len(b)).Info("message received")
					if err != nil {
						continue
					}
					msgC <- b
				}
			}
		}(c)
	}
}

func (d *Deamon) getClient(ctx context.Context) (*Client, error) {
	select {

	case c := <-d.clients:
		return c, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (d *Deamon) putClient(ctx context.Context, c *Client) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case d.clients <- c:
		return nil
	}
}
func (d *Deamon) recycleClient(ctx context.Context, c *Client) *Client {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				newC, err := NewClient(d.addr, d.username, d.password)
				if err != nil {
					logrus.Errorf("cannot create client: %v", err)
					continue
				}
				d.clients <- newC
			}
		}
	}()

	c.Close()
	return <-d.clients
}

func (d *Deamon) startSending(ctx context.Context, msgC chan []byte) {
	nclients := 50
	emailClient := make(chan *Client, nclients)
	clientErrC := make(chan error, nclients)

	go func() {
		for i := 0; i < nclients; i++ {
			client, err := NewClient(d.addr, d.username, d.password)
			if err != nil {
				logrus.Errorf("cannot create client: %v", err)
				continue
			}
			logrus.Debug("client created")
			emailClient <- client
		}
	}()

	go func() {
		logrus.Debug("client creation routine started")
		for {
			select {
			case <-ctx.Done():
				return
			default:
				<-clientErrC
				logrus.Info("client failed. creating a new one")
				go func() {
					created := false
					for !created {
						client, err := NewClient(d.addr, d.username, d.password)
						if err != nil {
							logrus.Errorf("cannot create client: %v", err)
							continue
						}
						emailClient <- client
						created = true
						logrus.Info("client created")
					}
				}()
			}
		}
	}()

	go func() {
		logrus.Debug("email sending routine started")
		for {
			select {
			case <-ctx.Done():
				return
			default:
				b := <-msgC
				logrus.WithField("msg_size", len(b)).Debug("message about to be send")
				go func(b []byte) {
					e, err := Unmarshal(b)
					if err != nil {
						logrus.Errorf("cannot parse email: %v", err)
						return
					}
					sent := false
					for i := 0; i < 5 && !sent; i++ {
						client := <-emailClient
						err := client.Send(e)
						if err != nil {
							logrus.WithField("email_id", e.ID()).Errorf("failed to send email: %v", err)
							e.AddStatusEvent(MakeStatusEvent(FailedAttemptToSend, time.Now()))
							client.Close()
							clientErrC <- err
							continue
						}
						logrus.WithField("email_id", e.ID()).Info("Email send")
						sent = true
						emailClient <- client
					}
					if sent {
						e.AddStatusEvent(MakeStatusEvent(SentSuccessfully, time.Now()))
					} else {
						e.AddStatusEvent(MakeStatusEvent(Dead, time.Now()))
					}
				}(b)
			}
		}
	}()
}
