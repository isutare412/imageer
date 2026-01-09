package sqs

import "time"

type ImageUploadListenerConfig struct {
	QueueURL           string
	BatchCount         int32
	PollingWaitTimeout time.Duration
	VisibilityTimeout  time.Duration
	BatchHandleTimeout time.Duration
	HandleTimeout      time.Duration
}
