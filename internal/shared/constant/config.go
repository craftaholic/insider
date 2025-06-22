package constant

const (
	DefaultTimeout        = 30
	WriteTimeout          = 30
	IdleTimeout           = 120
	DefaultContextTimeOut = 30
	CorsMaxAge            = 30

	RestExponentialBackOffScale = 2
	RestMaxRetry                = 2
	RestRetryWaitTime           = 2
	RestMaxWaitTime             = 10
	RestTimeOut                 = 30

	WorkerDefaultChanBuffer = 100
	WorkerDefaultCount      = 5

	ProducerDefaultCronDuration = 30
	ProducerDefaultBatchNumber  = 2

	WebhookDefaultTimeout = 30
)
