package main

import (
	"log"

	rabbitmq "github.com/iamgoangle/go-advance-rabbitmq/pkg/rabbitmq"
	middlewares "github.com/iamgoangle/go-advance-rabbitmq/pkg/rabbitmq/middlewares"
)

func main() {
	conf := rabbitmq.Config{
		Host:     "localhost",
		Port:     5672,
		Username: "admin",
		Password: "1234",
		Vhost:    "/",
	}
	connection, err := rabbitmq.NewAMQPConnection(conf)
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
	err = producer.Publish([]byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`))
	if err != nil {
		log.Println("unable to publish body")
	}
}
