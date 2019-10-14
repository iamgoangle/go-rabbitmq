package rabbitmq

import (
	"github.com/pkg/errors"

	"github.com/streadway/amqp"
)

// Producer handles interface to publish the message to broker and queue
type Producer interface {
	// UseWithConfig configs producer
	UseWithConfig(config ...ProducerConfigHandler) error

	// Publish send the message to AMQP broker
	Publish(body []byte, config ...PublishConfigHandler) error
}

// Produce represents produce object
type Produce struct {
	exchange, key, kind  string
	mandatory, immediate bool
	ch                   *amqp.Channel
}

// ProducerConfigHandler handles optinal parameter as a function
// See https://godoc.org/github.com/streadway/amqp#Channel.Publish
type ProducerConfigHandler func(*Produce) error

// PublishConfigHandler handles pass the publish config func
// See https://godoc.org/github.com/streadway/amqp#Publishing
type PublishConfigHandler func(*amqp.Publishing) error

// NewProducer instance the new producer
func NewProducer(exName, routingKey, kind string, ch *amqp.Channel) Producer {
	return &Produce{
		exchange: exName,
		key:      routingKey,
		kind:     kind,
		ch:       ch,
	}
}

func (p *Produce) UseWithConfig(configs ...ProducerConfigHandler) error {
	for _, config := range configs {
		err := config(p)
		if err != nil {
			return errors.Wrap(err, "unable to apply config")
		}
	}

	return nil
}

func (p *Produce) Publish(body []byte, config ...PublishConfigHandler) error {
	msg := amqp.Publishing{
		ContentType:  "text/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	}

	return p.ch.Publish(p.exchange, p.key, p.mandatory, p.immediate, msg)
}
