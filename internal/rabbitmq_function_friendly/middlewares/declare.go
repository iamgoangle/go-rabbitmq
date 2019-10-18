package rabbitmq

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	rabbitmq "github.com/iamgoangle/go-advance-rabbitmq/internal/rabbitmq_function_friendly"
)

const (
	ExchangeDirect  = amqp.ExchangeDirect
	ExchangeFanout  = amqp.ExchangeFanout
	ExchangeTopic   = amqp.ExchangeTopic
	ExchangeHeaders = amqp.ExchangeHeaders
)

// ExchangeDeclare return exchange declare handler function option
func ExchangeDeclare(name, kind string, args amqp.Table) rabbitmq.HandlerFunc {
	durable := false
	autoDelete := false
	internal := false
	noWait := false

	// inject your business logic here

	return func(c rabbitmq.Connection) error {
		err := c.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
		if err != nil {
			return errors.Wrap(err, "unable to declare exchange")
		}

		return nil
	}
}

// QueueDeclare return queue declare handler function option
func QueueDeclare(name string, args amqp.Table) rabbitmq.HandlerFunc {
	durable := false
	autoDelete := false
	exclusive := false
	noWait := false

	// inject your business logic here

	return func(c rabbitmq.Connection) error {
		err := c.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
		if err != nil {
			return errors.Wrap(err, "unable to declare queue")
		}

		return nil
	}
}
