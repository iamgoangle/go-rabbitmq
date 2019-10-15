package main

import (
	"log"

	rabbitmq "github.com/iamgoangle/go-advance-rabbitmq/internal/rabbitmq_function_friendly"
)

func main() {
	connection, err := rabbitmq.NewAMQPConnection("amqp://admin:1234@localhost:5672/")
	if err != nil {
		log.Panic(err)
	}

	connection.Use(rabbitmq.ExchangeDeclare("exchange_test", rabbitmq.ExchangeDirect, nil))
	connection.Use(rabbitmq.QueueDeclare("test", nil))
	connection.Use(rabbitmq.QueueBind("test", "", "exchange_test", false, nil))

	producer := rabbitmq.NewProducer("exchange_test", "", "", connection.Services())
	err = producer.Publish([]byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`), nil)
	if err != nil {
		log.Println("unable to publish body")
	}
}
