package rabbitmq

type Consumer interface {
	Use()
	UseWithDLQ()
}
