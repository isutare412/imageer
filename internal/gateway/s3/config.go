package s3

import "time"

type PresignerConfig struct {
	Bucket string
	Expiry time.Duration
}
