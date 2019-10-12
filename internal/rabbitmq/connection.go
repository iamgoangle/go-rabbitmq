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
}

type connection struct {
	*amqp.Connection
	*amqp.Channel
}

type Exchange struct {
	name, kind                            string
	durable, autoDelete, internal, noWait bool
	args                                  *amqp.Table
}

// NewAMQPConnection creates amqp connection and channel
func NewAMQPConnection(host string) (Connection, error) {
	if len(host) == 0 {
		return nil, errors.New("fail to create channel due to missing host specified")
	}

	conn, err := amqp.Dial(host)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create new connection")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create channel")
	}

	log.Println("created new amqp connection and channel")

	return &connection{
		Connection: conn,
		Channel:    ch,
	}, nil
}

func (c *connection) Close() {
	c.Close()
}

func (c *connection) CloseChannel() {
	c.CloseChannel()
}

func (c *connection) ExchangeDeclare(name, kind string, configs ...Exchange) error {
	err := c.Channel.ExchangeDeclare(name, kind)
	if err != nil {
		return errors.Wrap(err, "unable to declare exchange")
	}

	return nil
}
