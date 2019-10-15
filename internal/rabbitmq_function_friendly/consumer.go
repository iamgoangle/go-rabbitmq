package rabbitmq

import (
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Consumer interface {
	Use(handler ConsumerHandler)

	// WithConfigs config consumer
	// See https://godoc.org/github.com/streadway/amqp#Channel.Consume
	WithConfigs(configs ...ConsumerConfigHandler)

	// WithDeadLetterQueue defines dead-letter-queue with friendly config
	WithDeadLetterQueue()

	Consume() error

	ConsumeWithRetry()
}

type ConsumerConfigHandler func(*Consume) error

type Consume struct {
	ch *amqp.Channel

	queueName    string
	consumerName string
	autoAck      bool
	exclusive    bool
	noLocal      bool
	noWait       bool
	args         amqp.Table

	msg      chan *amqp.Delivery
	handlers []ConsumerHandler
}

type ConsumerHandler interface {
	Do(msg []byte) error

	Fallback(err error)
}

// NewConsumer creates an instance the consumer object
// qName specific queue name you want to consume
// cName specific consumer name
func NewConsumer(qName, cName string, ch *amqp.Channel) Consumer {
	return &Consume{
		queueName:    qName,
		consumerName: cName,
		ch:           ch,
	}
}

func (c *Consume) WithConfigs(configs ...ConsumerConfigHandler) {
	for _, config := range configs {
		err := config(c)
		if err != nil {
			log.Panic(errors.Wrap(err, FailedToApplyConsumerConfigFunc))
		}
	}
}

func (c *Consume) WithDeadLetterQueue() {

}

func (c *Consume) Use(handler ConsumerHandler) {
	if handler == nil {
		log.Panic(FailedToAppledHandlerFunc)
	}

	c.handlers = append(c.handlers, handler)
}

func (c *Consume) Consume() error {
	msgs, err := c.ch.Consume(
		c.queueName,    // queue
		c.consumerName, // consumer
		c.autoAck,      // auto ack
		c.exclusive,    // exclusive
		c.noLocal,      // no local
		c.noWait,       // no wait
		c.args,         // args
	)

	if err != nil {
		return errors.Wrap(err, FailedToRegisterConsumer)
	}

	go func() {
		for m := range msgs {
			for _, h := range c.handlers {
				err := h.Do(m.Body)
				if err != nil {
					h.Fallback(err)
					break
				}
			}
		}
	}()

	c.ch.Close()

	return nil
}

func (c *Consume) ConsumeWithRetry() {

}
