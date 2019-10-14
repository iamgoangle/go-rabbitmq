package main

import (
	"log"

	rabbitmq "github.com/iamgoangle/go-advance-rabbitmq/internal/rabbitmq_function_friendly"
)

func main() {
	connection, err := rabbitmq.NewAMQPConnection("test")
	if err != nil {
		log.Panic(err)
	}

	connection.Declare(rabbitmq.Exchange("exchange_test", "direct", nil))
	connection.Declare(rabbitmq.Queue("test", nil))

	connection.
}
