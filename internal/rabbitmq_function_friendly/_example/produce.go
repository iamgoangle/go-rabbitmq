package main

import (
	"log"

	rabbitmq "github.com/iamgoangle/go-advance-rabbitmq/internal/rabbitmq_function_friendly"
	middlewares "github.com/iamgoangle/go-advance-rabbitmq/internal/rabbitmq_function_friendly/middlewares"
)

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

	producer := rabbitmq.NewProducer("exchange_test", "", "", connection)
	err = producer.Publish([]byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`), nil)
	if err != nil {
		log.Println("unable to publish body")
	}
}
