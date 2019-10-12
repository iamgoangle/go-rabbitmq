package rabbitmq

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Producer interface {
	UseExchangeWithConfig(...ExchangeConfigHandler) error
}

type producer struct {
	*Exchange
}

type Exchange struct {
	name, kind                            string
	durable, autoDelete, internal, noWait bool
	args                                  *amqp.Table
}

type ExchangeConfigHandler func(*Exchange) error

// NewProducer instance the new producer
func NewProducer(exName, exKind string, conn Connection) Producer {
	return &producer{
		Exchange: &Exchange{
			name: exName,
			kind: exKind,
		},
	}
}

func (p *producer) UseExchangeWithConfig(configs ...ExchangeConfigHandler) error {
	for _, config := range configs {
		err := config(p.Exchange)
		if err != nil {
			return errors.Wrap(err, "unable to apply exchange with config")
		}
	}

	return nil
}
