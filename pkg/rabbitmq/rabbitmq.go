package rabbitmq

//go:generate mockgen -source=./rabbitmq.go -destination=./mocks/rabbitmq_mock.go -package=mocks github.com/iamgoangle/go-advance-rabbitmq/pkg/rabbitmq Connection

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	delay = 1
)

// Config is RabbitMQ connection
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Vhost    string
}

// Connection represent interface RabbitMQ method
type Connection interface {
	// Use push middleware function apply to connection
	// middleware can be config function
	Use(handler HandlerFunc) error

	// Close entire amqp connection
	Close()

	// CloseChannel closes the amqp channel
	CloseChannel()

	// ApplyUse applies use as soon as possible
	ApplyUse(handler ...HandlerFunc) error

	// Declare declares an exchange or queue
	Declare

	// Bind binds the queue with exchange
	Bind

	// Channel represents an AMQP channel.
	// Used as a context for valid message exchange.
	// Errors on methods with this Channel as a receiver means this channel
	// should be discarded and a new channel established.
	Channel

	// Run apply middlewares
	Run() error
}

type Declare interface {
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
}

type Bind interface {
	// QueueBind binds an exchange to a queue so that publishings to the exchange will
	// be routed to the queue when the publishing routing key matches the binding
	// routing key.
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
}

type Channel interface {
	// Publish sends a Publishing from the client to an exchange on the server
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error

	// Consume immediately starts delivering queued messages
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}

// HandlerFunc is a handler function for decorate exchange and queue
// that allows the client can be modify thier content as a closure
// and return middlewares function
type HandlerFunc func(c Connection) error

type connection struct {
	*amqp.Connection
	Channel *channel

	middlewares []HandlerFunc
}

type channel struct {
	*amqp.Channel
	closed int32
}

// NewAMQPConnection creates amqp connection and channel
func NewAMQPConnection(c Config) (Connection, error) {
	if c.Host == "" {
		return nil, errors.New("fail to create channel due to missing host specified")
	}

	conf := amqp.URI{
		Scheme:   "amqp",
		Host:     c.Host,
		Port:     c.Port,
		Username: c.Username,
		Password: c.Password,
		Vhost:    c.Vhost,
	}.String()

	conn, err := amqp.Dial(conf)
	if err != nil {
		return nil, errors.Wrap(err, FailedToCreateNewConnection)
	}

	reconnectable := &connection{
		Connection: conn,
	}

	go func() {
		rabbitError := make(chan *amqp.Error)

		for {
			reason, ok := <-reconnectable.Connection.NotifyClose(rabbitError)
			// exit this goroutine if closed by developer
			if !ok {
				log.Println("connection closed")
				break
			}
			log.Printf("connection closed, reason: %v \n", reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for reconnect
				time.Sleep(delay * time.Second)

				conn, err := amqp.Dial(conf)
				if err == nil {
					reconnectable.Connection = conn
					log.Println("reconnect success")
					break
				}

				log.Printf("reconnect failed, err: %v \n", err)
			}
		}
	}()

	ch, err := reconnectable.channel()
	if err != nil {
		return nil, errors.Wrap(err, FailedToCreateNewChannel)
	}

	reconnectable = &connection{
		Channel: ch,
	}

	log.Printf("created new amqp connection and channel %v", conf)

	return &connection{
		Connection: conn,
		Channel:    ch,
	}, nil
}

func (c *connection) channel() (*channel, error) {
	ch, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	channel := &channel{
		Channel: ch,
	}

	go func() {
		for {
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok || channel.IsClosed() {
				log.Println("channel closed")
				channel.Close() // close again, ensure closed flag set when connection closed
				break
			}
			log.Printf("channel closed, reason: %v \n", reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for connection reconnect
				time.Sleep(delay * time.Second)

				ch, err := c.Connection.Channel()
				if err == nil {
					log.Println("channel recreate success")
					channel.Channel = ch
					break
				}

				log.Printf("channel recreate failed, err: %v \n", err)
			}
		}

	}()

	return channel, nil
}

// IsClosed indicate closed by developer
func (ch *channel) IsClosed() bool {
	return (atomic.LoadInt32(&ch.closed) == 1)
}

// Close ensure closed flag set
func (ch *channel) Close() error {
	if ch.IsClosed() {
		return amqp.ErrClosed
	}

	atomic.StoreInt32(&ch.closed, 1)

	return ch.Channel.Close()
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

func (c *connection) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return c.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (c *connection) Close() {
	if err := c.Connection.Close(); err != nil {
		log.Panic("unable to close connection")
	}

	log.Println("Closed AMQP Channel")
}

func (c *connection) CloseChannel() {
	if err := c.Channel.Close(); err != nil {
		log.Panic("unable to close channel")
	}

	log.Println("Closed AMQP Connection")
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
