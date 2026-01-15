package image

import "time"

type Config struct {
	CDNDomain              string
	S3KeyPrefix            string
	ProcessDoneWaitTimeout time.Duration
}

type CloserConfig struct {
	CheckInterval  time.Duration
	CloseThreshold time.Duration
}
