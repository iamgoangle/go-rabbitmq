package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"

	rabbitmq "github.com/iamgoangle/go-advance-rabbitmq/pkg/rabbitmq"
	middlewares "github.com/iamgoangle/go-advance-rabbitmq/pkg/rabbitmq/middlewares"
)

type handler struct {
}

func main() {
	connection, err := rabbitmq.NewAMQPConnection("amqp://admin:1234@localhost:5672/")
	if err != nil {
		log.Panic(err)
	}

	connection.Use(middlewares.ExchangeDeclare("exchange_test", middlewares.ExchangeDirect, nil))
	connection.Use(middlewares.QueueDeclare("test", nil))
	connection.Use(middlewares.QueueBind("test", "", "exchange_test", false, nil))
	if err := connection.Run(); err != nil {
		log.Panic(err)
	}

	consumer := rabbitmq.NewConsumer("test", "consumer_name", connection)
	consumer.Use(newConsumerHandler())

	if err := consumer.Consume(); err != nil {
		log.Panic(err)
	}
}

func newConsumerHandler() rabbitmq.ConsumerHandler {
	return &handler{}
}

func (h *handler) Do(msg []byte) error {
	fmt.Println(string(msg))

	return nil
}

func (h *handler) OnSuccess(m amqp.Delivery) error {
	log.Println("consume item success")

	return nil
}

func (h *handler) OnError(m amqp.Delivery, err error) {

}
