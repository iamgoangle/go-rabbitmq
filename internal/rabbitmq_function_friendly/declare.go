package rabbitmq

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	ExchangeDirect  = "direct"
	ExchangeFanout  = "fanout"
	ExchangeTopic   = "topic"
	ExchangeHeaders = "headers"
)

// ExchangeDeclare return exchange declare handler function option
func ExchangeDeclare(name, kind string, args amqp.Table) HandlerFunc {
	durable := false
	autoDelete := false
	internal := false
	noWait := false

	return func(ch *amqp.Channel) error {
		err := ch.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
		if err != nil {
			return errors.Wrap(err, "unable to declare exchange")
		}

		return nil
	}
}

// QueueDeclare return queue declare handler function option
func QueueDeclare(name string, args amqp.Table) HandlerFunc {
	durable := false
	autoDelete := false
	exclusive := false
	noWait := false

	return func(ch *amqp.Channel) error {
		_, err := ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
		if err != nil {
			return errors.Wrap(err, "unable to declare queue")
		}

		return nil
	}
}
