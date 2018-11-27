package messaging

import (
	"fmt"
	"net"
	"time"

	"github.com/streadway/amqp"
)

// RabbitMqConnection is an implementation of Publisher for rabbit mq.
type RabbitMqConnection struct {
	conn    *amqp.Connection
	connStr string
	queue   string
}

// NewRabbitMqConnection creates a new instance of RabbitMqConnection
func NewRabbitMqConnection(connStr, queue string) (*RabbitMqConnection, error) {
	conn, err := amqp.DialConfig(connStr, amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 5*time.Second)
		},
	})
	if err != nil {
		return nil, err
	}
	rmc := &RabbitMqConnection{
		conn:    conn,
		connStr: connStr,
		queue:   queue,
	}

	if err := rmc.queueDeclare(); err != nil {
		return nil, fmt.Errorf("cannot declare queue %s: %v", rmc.queue, err)
	}
	return rmc, nil
}

// Publish to rabbit mq queue
func (c *RabbitMqConnection) Publish(b []byte) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("cannot open channel: %v", err)
	}

	defer ch.Close()
	return ch.Publish("", c.queue, false, false, amqp.Publishing{Body: b})
}

// Close the rabbit mq connection
func (c *RabbitMqConnection) Close() error {
	return c.conn.Close()
}

func (c *RabbitMqConnection) queueDeclare() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create a channel: %v", err)
	}
	_, err = ch.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}
