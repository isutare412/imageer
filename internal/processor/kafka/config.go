package kafka

import "time"

type ClientConfig struct {
	Addrs         []string
	User          string
	Password      string
	ConsumerGroup string
	ConsumeTopics []string
}

type ImageProcessRequestHandlerConfig struct {
	RetryTopic      string
	HandleTimeout   time.Duration
	MaxRetryAttempt int
	RetryBaseDelay  time.Duration
}

type ImageProcessResultQueueConfig struct {
	Topic string
}
