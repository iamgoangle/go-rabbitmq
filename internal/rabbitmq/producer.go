package rabbitmq

import (
	"github.com/pkg/errors"

	"github.com/streadway/amqp"
)

type Producer interface {
	UseWithConfig(config ...ProducerConfigHandler) error
	Publish(body []byte, config ...PublishConfigHandler) error
}

type Produce struct {
	exchange, key, kind  string
	mandatory, immediate bool
	con                  Connection
}

type ProducerConfigHandler func(*Produce) error

type PublishConfigHandler func(*amqp.Publishing) error

// NewProducer instance the new producer
func NewProducer(exName, routingKey, kind string, con Connection) Producer {
	return &Produce{
		exchange: exName,
		key:      routingKey,
		kind:     kind,
		con:      con,
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
		Body: body,
	}

	return p.con.Do().Publish(p.exchange, p.key, p.mandatory, p.immediate, msg)
}
