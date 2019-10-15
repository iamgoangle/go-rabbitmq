package rabbitmq

import (
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Consumer interface {
	Use()

	// UseWithConfigs config consumer
	// See https://godoc.org/github.com/streadway/amqp#Channel.Consume
	UseWithConfigs(configs ...ConsumerConfigHandler)

	UseDeadLetterQueue()
}

type ConsumerConfigHandler func(*Consume) error

type Consume struct {
	queueName    string
	consumerName string
	autoAck      bool
	exclusive    bool
	noLocal      bool
	noWait       bool
	args         *amqp.Table
}

// NewConsumer creates an instance the consumer object
// qName specific queue name you want to consume
// cName specific consumer name
func NewConsumer(qName, cName string, ch *amqp.Channel) Consumer {
	return &Consume{
		queueName:    qName,
		consumerName: cName,
	}
}

func (c *Consume) Use() {

}

func (c *Consume) UseWithConfigs(configs ...ConsumerConfigHandler) {
	for _, config := range configs {
		err := config(c)
		if err != nil {
			log.Panic(errors.Wrap(err, FailedToApplyConsumerConfigFunc))
		}
	}
}

func (c *Consume) UseDeadLetterQueue() {

}
