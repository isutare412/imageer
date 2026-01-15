package image

import "time"

type Config struct {
	CDNDomain              string
	S3KeyPrefix            string
	ProcessDoneWaitTimeout time.Duration
}
