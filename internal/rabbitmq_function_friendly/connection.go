package rabbitmq

import (
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type HandlerFunc func(ch *amqp.Channel) error

// Connection represent interface amqp connection
type Connection interface {
	// Close entire amqp connection
	Close()

	// CloseChannel closes the amqp channel
	CloseChannel()

	// Do handles channel amqp event
	Do() *amqp.Channel

	// Use applys use handler
	// exchange to define amqp routing
	// queue to define queue
	// Use(handlers ...HandlerFunc) error
	Use(handler HandlerFunc) error

	// ApplyUse applies use as soon as possible
	ApplyUse(handler ...HandlerFunc) error
}

type connection struct {
	*amqp.Connection
	*amqp.Channel

	middlewares []HandlerFunc
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

func (c *connection) Use(handler HandlerFunc) error {
	if handler == nil {
		log.Panic("unable to apply Uses")
	}

	c.Uses = append(c.Uses, handler)

	return nil
}

func (c *connection) ApplyUse(handlers ...HandlerFunc) error {
	if handlers == nil {
		log.Panic("unable to apply Uses")
	}

	for _, h := range handlers {
		err := h(c.Channel)
		if err != nil {
			return errors.Wrap(err, "unable appled Use handler")
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
	for _, handler := range c.middlewares {
		err := handler(c.Channel)
		if err != nil {
			log.Panic("unable applies handler", err.Error())
		}
	}

	return c.Channel
}
