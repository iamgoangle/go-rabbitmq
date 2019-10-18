package rabbitmq

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	rabbitmq "github.com/iamgoangle/go-advance-rabbitmq/internal/rabbitmq_function_friendly"
)

// QueueBind binding queue with exchange
func QueueBind(name, key, exchange string, noWait bool, args amqp.Table) rabbitmq.HandlerFunc {
	// inject your business logic here

	return func(c rabbitmq.Connection) error {
		err := c.QueueBind(name, key, exchange, noWait, args)
		if err != nil {
			return errors.Wrap(err, "unable to binding queue")
		}

		return nil
	}
}
