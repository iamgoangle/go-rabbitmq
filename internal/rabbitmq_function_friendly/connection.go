package rabbitmq

import (
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// HandlerFunc is handler function type for inject a logic
// for implment amqp method

// type HandlerFunc func(ch *amqp.Channel) error
type HandlerFunc func(c Connection) error

// Connection represent interface amqp connection
type Connection interface {
	// Close entire amqp connection
	Close()

	// CloseChannel closes the amqp channel
	CloseChannel()

	// Use applys use handler
	// exchange to define amqp routing
	// queue to define queue
	// Use(handlers ...HandlerFunc) error
	Use(handler HandlerFunc) error

	// ApplyUse applies use as soon as possible
	ApplyUse(handler ...HandlerFunc) error

	Declare
	Bind
	Channel

	Run() error
}

// Declare handler amqp channel declare
type Declare interface {
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error

	QueueDeclare(name string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
}

// Bind handle amqp channel binding
type Bind interface {
	// QueueBind binds an exchange to a queue so that publishings to the exchange will
	// be routed to the queue when the publishing routing key matches the binding
	// routing key.
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
}

type Channel interface {
	// Publish sends a Publishing from the client to an exchange on the server.
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
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
		return nil, errors.Wrap(err, FailedToCreateNewConnection)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, FailedToCreateNewChannel)
	}

	log.Println("created new amqp connection and channel")

	return &connection{
		Connection: conn,
		Channel:    ch,
	}, nil
}

func (c *connection) Use(handler HandlerFunc) error {
	if handler == nil {
		log.Panic(FailedToAppledHandlerFunc)
	}

	c.middlewares = append(c.middlewares, handler)

	return nil
}

func (c *connection) ApplyUse(handlers ...HandlerFunc) error {
	if handlers == nil {
		log.Panic(FailedToAppledHandlerFunc)
	}

	for _, h := range handlers {
		err := h(c)
		if err != nil {
			return errors.Wrap(err, FailedToAppledHandlerFunc)
		}
	}

	return nil
}

func (c *connection) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	err := c.Channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
	if err != nil {
		return errors.Wrap(err, "unable to declare exchange")
	}

	return nil
}

func (c *connection) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) error {
	_, err := c.Channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
	if err != nil {
		return errors.Wrap(err, "unable to declare queue")
	}

	return nil
}

func (c *connection) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	err := c.Channel.QueueBind(name, key, exchange, noWait, args)
	if err != nil {
		return errors.Wrap(err, "unable to binding queue")
	}

	return nil
}

func (c *connection) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return c.Channel.Publish(exchange, key, mandatory, immediate, msg)
}

func (c *connection) Close() {
	c.Close()
}

func (c *connection) CloseChannel() {
	c.CloseChannel()
}

func (c *connection) Run() error {
	for _, handler := range c.middlewares {
		err := handler(c)
		if err != nil {
			errors.Wrap(err, FaiiledToRun)
		}
	}

	return nil
}
