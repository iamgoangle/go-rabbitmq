package rabbitmq

const (
	FailedToCreateNewConnection = "unable to create new connection"
	FailedToCreateNewChannel    = "unable to create channel"
	FailedToAppledHandlerFunc   = "unable to apply handler func, please check your function"

	FailedToApplyProducerConfigFunc        = "unable to apply config"
	FailedToApplyPublishPropertyConfigFunc = "unable to apply publish property config func"

	MissingArgumentTTL         = "missing argument `ttl`"
	MissingArgumentContentType = "missing argument `cType`"
	FailedToSetConfigPersist   = "failed to set config `persist`"
	FailedToSetConfigPriority  = "fail to set config `prority queue`"
)
