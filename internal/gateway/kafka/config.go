package kafka

import "time"

type ClientConfig struct {
	Addrs         []string
	User          string
	Password      string
	ConsumerGroup string
	ConsumeTopics []string
}

type ImageProcessRequestQueueConfig struct {
	Topic string
}

type ImageProcessResultHandlerConfig struct {
	RetryTopic      string
	HandleTimeout   time.Duration
	MaxRetryAttempt int
	RetryBaseDelay  time.Duration
}

type ImageS3DeleteRequestQueueConfig struct {
	Topic string
}

type ImageS3DeleteRequestHandlerConfig struct {
	RetryTopic      string
	HandleTimeout   time.Duration
	MaxRetryAttempt int
	RetryBaseDelay  time.Duration
}
