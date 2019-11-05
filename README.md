# GO AMQP Client Wrapper

## About

Implementation is probably straight-forward of the project that aim to follow chain-of-responsiblity and [functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) pattern to solve the complexity and make abstraction layer to a user that use the this library.

### Features

- [x] Middlewares for **exchange**, **queue** and **queue bind**.
- [x] Functional options for producer publishing the message.
- [x] Re-establish connection and channel to Publisher and Consumer.
- [ ] Consumer retry queue middleware.
- [ ] Consumer retry queue middleware with backoff config.
- [x] Closed **connection** and **channel** gracefully.

### Producer

```shell
import (
	"log"

	rabbitmq "github.com/iamgoangle/go-rabbitmq/pkg/rabbitmq"
	middlewares "github.com/iamgoangle/go-rabbitmq/pkg/rabbitmq/middlewares"
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
```

### Basic Consumer

```sh
import (
	"fmt"
	"log"

	rabbitmq "github.com/iamgoangle/go-rabbitmq/pkg/rabbitmq"
	middlewares "github.com/iamgoangle/go-rabbitmq/pkg/rabbitmq/middlewares"
)

type handler struct {
}

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
```

