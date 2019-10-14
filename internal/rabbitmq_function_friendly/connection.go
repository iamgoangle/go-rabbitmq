package rabbitmq

import (
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type DeclareHandler func(ch *amqp.Channel) error

// Connection represent interface amqp connection
type Connection interface {
	// Close entire amqp connection
	Close()

	// CloseChannel closes the amqp channel
	CloseChannel()

	// Do handles channel amqp event
	Do() *amqp.Channel

	// Declare applys declare handler
	// exchange to define amqp routing
	// queue to define queue
	// Declare(handlers ...DeclareHandler) error
	Declare(handler DeclareHandler) error

	// ApplyDeclare applies declare as soon as possible
	ApplyDeclare(handler ...DeclareHandler) error
}

type connection struct {
	*amqp.Connection
	*amqp.Channel

	declares []DeclareHandler
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

func (c *connection) Declare(handler DeclareHandler) error {
	if handler == nil {
		log.Panic("unable to apply declares")
	}

	c.declares = append(c.declares, handler)

	return nil
}

func (c *connection) ApplyDeclare(handlers ...DeclareHandler) error {
	if handlers == nil {
		log.Panic("unable to apply declares")
	}

	for _, h := range handlers {
		err := h(c.Channel)
		if err != nil {
			return errors.Wrap(err, "unable appled declare handler")
		}
	}

	return nil
}

func (c *connection) Close() {
	c.Close()
}

func (c *connection) CloseChannel() {
	c.CloseChannel()
}

func (c *connection) Do() *amqp.Channel {
	for _, h := range c.declares {
		err := h(c.Channel)
		if err != nil {
			log.Panic("unable applies declare handler", err.Error())
		}
	}

	return c.Channel
}
