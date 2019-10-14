package main

import (
	"fmt"

	"github.com/iamgoangle/go-advance-rabbitmq/internal/rabbitmq"
)

func main() {
	connection := rabbitmq.NewAMQPConnection("test").
		ExchangeDeclare("exchange_test", "direct", nil).
		QueueDeclare("test", nil)

	// do produce or consume...
	fmt.Println(connection)
}
