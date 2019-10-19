package rabbitmq

const (
	FailedToCreateNewConnection = "unable to create new connection"
	FailedToCreateNewChannel    = "unable to create channel"
	FailedToAppledHandlerFunc   = "unable to apply handler func, please check your function"
	FaiiledToRun                = "unable to execute middleware"

	FailedToApplyProducerConfigFunc        = "unable to apply producer config"
	FailedToApplyPublishPropertyConfigFunc = "unable to apply publish property config func"

	MissingArgumentTTL         = "missing argument `ttl`"
	MissingArgumentContentType = "missing argument `cType`"
	FailedToSetConfigPersist   = "failed to set config `persist`"
	FailedToSetConfigPriority  = "fail to set config `prority queue`"

	FailedToApplyConsumerConfigFunc = "unable to apply consumer config"
	FailedToRegisterConsumer        = "failed to register consumer"

	TTLShouldBeGTZero               = "ttl must be greater than zero"
	FailedToExecuteConsumerHandlers = "failed to execute consumer handlers"
	FailedToConsumeItem             = "failed to consume message"
)

const (
	ConsumerRegistered = "registered new consumer"
)
