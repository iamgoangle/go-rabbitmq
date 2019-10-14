package rabbitmq

import (
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Connection represent interface amqp connection
type Connection interface {
	// Close entire amqp connection
	Close()

	// CloseChannel closes the amqp channel
	CloseChannel()

	// Do handles channel amqp event
	Do() *amqp.Channel

	// ExchangeDeclare creates an exchange
	ExchangeDeclare(name, kind string, args amqp.Table) Connection

	// QueueDeclare creates a queue
	QueueDeclare(name string, args amqp.Table) Connection
}

type connection struct {
	*amqp.Connection
	*amqp.Channel
}

// NewAMQPConnection creates amqp connection and channel
func NewAMQPConnection(host string) Connection {
	if len(host) == 0 {
		log.Panic(errors.New("fail to create channel due to missing host specified"))
	}

	conn, err := amqp.Dial(host)
	if err != nil {
		log.Panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Panic(err)
	}

	log.Println("created new amqp connection and channel")

	return &connection{
		Connection: conn,
		Channel:    ch,
	}
}

func (c *connection) Close() {
	c.Close()
}

func (c *connection) CloseChannel() {
	c.CloseChannel()
}

func (c *connection) ExchangeDeclare(name, kind string, args amqp.Table) Connection {
	durable := false
	autoDelete := false
	internal := false
	noWait := false

	err := c.Channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
	if err != nil {
		log.Panic("unable to create exchange", err.Error())
	}

	return c
}

func (c *connection) Do() *amqp.Channel {
	return c.Channel
}

func (c *connection) QueueDeclare(name string, args amqp.Table) Connection {
	durable := false
	autoDelete := false
	exclusive := false
	noWait := false

	_, err := c.Channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
	if err != nil {
		log.Panic("unable to create queue", err.Error())
	}

	return c
}
