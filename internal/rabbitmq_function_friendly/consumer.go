package rabbitmq

type Consumer interface {
	Use()
	UseDeadLetterQueue()
}
