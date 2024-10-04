package rabbitmq

// https://medium.com/@dhanushgopinath/automatically-recovering-rabbitmq-connections-in-go-applications-7795a605ca59

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

// MessageBody is the struct for the body passed in the AMQP message. The type will be set on the Request header
type MessageBody struct {
	Data []byte
	Type string
}

// Message is the amqp request to publish
type Message struct {
	Queue         string
	ReplyTo       string
	ContentType   string
	CorrelationId string
	Priority      uint8
	Body          MessageBody
}

// Connection is the connection created
type Connection struct {
	cfg          amqp.URI
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchange     string
	exchangeType string
	queue        string
	routingKey   string
	qos          int
	err          chan error
}

// NewConnection returns the new connection object
func NewConnection(cfg amqp.URI, exchange, exchangeType, queue string, qos int) *Connection {
	c := &Connection{
		cfg:          cfg,
		exchange:     exchange,
		exchangeType: exchangeType,
		queue:        queue,
		qos:          qos,
		err:          make(chan error),
	}
	return c
}

func (c *Connection) connect() error {
	var err error
	c.conn, err = amqp.Dial(c.cfg.String())

	if err != nil {
		return fmt.Errorf("error in creating rabbitmq connection with %s : %s", c.cfg.String(), err.Error())
	}

	go func() {
		<-c.conn.NotifyClose(make(chan *amqp.Error)) //Listen to NotifyClose
		c.err <- errors.New("Connection Closed")
	}()

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	c.channel.Qos(c.qos, 0, false)

	if c.exchange != "" {
		if err := c.channel.ExchangeDeclare(
			c.exchange,     // name
			c.exchangeType, // type
			true,           // durable
			false,          // auto-deleted
			false,          // internal
			false,          // noWait
			nil,            // arguments
		); err != nil {
			return fmt.Errorf("error in Exchange Declare: %s", err)
		}
	}
	return nil
}

func (c *Connection) BindQueue() error {
	if c.queue != "" {
		_, err := c.channel.QueueDeclare(c.queue, true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("error in declaring the queue %s", err)
		}

		if c.exchange != "" {
			err := c.channel.QueueBind(c.queue, "", c.exchange, false, nil)
			if err != nil {
				return fmt.Errorf("queue  Bind error: %s", err)
			}
		}
	}

	return nil
}

// Reconnect reconnects the connection
func (c *Connection) Reconnect() error {
	if err := c.connect(); err != nil {
		return err
	}
	if err := c.BindQueue(); err != nil {
		return err
	}
	return nil
}
