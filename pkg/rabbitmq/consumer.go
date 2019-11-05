package rabbitmq

//go:generate mockgen -source=./consumer.go -destination=./mocks/consumer_mock.go -package=mocks github.com/iamgoangle/go-advance-rabbitmq/pkg/rabbitmq Consumer

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// ConsumerHandler defines the set of method for any consumer client
// provides Do() for client business logic
// and OnError() for handle in case of consume error
type ConsumerHandler interface {
	Do(msg []byte) error

	OnSuccess(m amqp.Delivery) error

	OnError(m amqp.Delivery, err error)
}

// Consumer interface
type Consumer interface {
	// WithConfigs config consumer
	// See https://godoc.org/github.com/streadway/amqp#Channel.Consume
	WithConfigs(configs ...ConsumerConfigHandler)

	// WithDeadLetterQueue defines dead-letter-queue with friendly config
	WithDeadLetterQueue(configs ...ConsumerConfigDLQHandler)

	// Use apply consumer handlers
	Use(handler ConsumerHandler)

	Consume() error

	ConsumeWithRetry()
}

// ConsumerConfigHandler is config function contain a implement method amqp.Channel
type ConsumerConfigHandler func(*Consume) error

type ConsumerConfigDLQHandler func(*ConsumerDLQ) error

// Consume type
type Consume struct {
	conn Connection

	queueName    string
	consumerName string
	autoAck      bool
	exclusive    bool
	noLocal      bool
	noWait       bool
	args         amqp.Table

	msg chan *amqp.Delivery

	handlers []ConsumerHandler

	requiredRetry bool

	*ConsumerDLQ
}

type ConsumerDLQ struct {
}

// NewConsumer creates an instance the consumer object
// qName specific queue name you want to consume
// cName specific consumer name
func NewConsumer(qName, cName string, conn Connection) Consumer {
	return &Consume{
		queueName:    qName,
		consumerName: cName,
		conn:         conn,
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

func (c *Consume) WithDeadLetterQueue(configs ...ConsumerConfigDLQHandler) {
	c.requiredRetry = true
}

func (c *Consume) Use(handler ConsumerHandler) {
	if handler == nil {
		log.Panic(FailedToAppledHandlerFunc)
	}

	c.handlers = append(c.handlers, handler)
}

func (c *Consume) Consume() error {
	done := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)

	defer close(done)
	defer close(sigs)

	msgs, err := c.conn.Consume(
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
		log.Printf("%s - %s\n", c.consumerName, ConsumerRegistered)

		for m := range msgs {
			if len(c.handlers) == 0 {
				log.Panic(FailedToExecuteConsumerHandlers)
			}

			for _, h := range c.handlers {
				err := h.Do(m.Body)
				if err != nil {
					h.OnError(m, err)
					break
				}

				if err := h.OnSuccess(m); err != nil {
					break
				}
			}

			m.Ack(false)
		}
	}()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("%v - recieved signals", sig)
		c.conn.Close()

		done <- true
	}()

	<-done

	log.Println("Closed rabbitmq consumer")

	return nil
}

func (c *Consume) ConsumeWithRetry() {

}
