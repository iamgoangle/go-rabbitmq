package rabbitmq

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Exchange return exchange declare handler function option
func Exchange(name, kind string, args amqp.Table) DeclareHandler {
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

// Queue return queue declare handler function option
func Queue(name string, args amqp.Table) DeclareHandler {
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
