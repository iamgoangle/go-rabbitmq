package rabbitmq

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// QueueBind binding queue with exchange
func QueueBind(name, key, exchange string, noWait bool, args amqp.Table) HandlerFunc {
	return func(ch *amqp.Channel) error {
		err := ch.QueueBind(name, key, exchange, noWait, args)
		if err != nil {
			return errors.Wrap(err, "unable to binding queue")
		}

		return nil
	}
}
