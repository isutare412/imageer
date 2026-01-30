package config

import (
	"time"

	"github.com/isutare412/imageer/pkg/log"
)

type Config struct {
	Log    LogConfig    `koanf:"log"`
	Web    WebConfig    `koanf:"web"`
	Valkey ValkeyConfig `koanf:"valkey"`
	AWS    AWSConfig    `koanf:"aws"`
}

type LogConfig struct {
	Format    log.Format `koanf:"format" validate:"validateFn=Validate"`
	Level     log.Level  `koanf:"level" validate:"validateFn=Validate"`
	AddSource bool       `koanf:"add-source"`
	Component string     `koanf:"component" validate:"required"`
}

type WebConfig struct {
	Port int `koanf:"port" validate:"required,gt=0,lte=65535"`
}

type ValkeyConfig struct {
	Addresses string `koanf:"addresses" validate:"required"`
	Username  string `koanf:"username" validate:"required"`
	Password  string `koanf:"password" validate:"required"`

	Streams struct {
		ImageProcessRequest struct {
			StreamKey string `koanf:"stream-key" validate:"required"`
			GroupName string `koanf:"group-name" validate:"required"`
			Handler   struct {
				Concurrency int           `koanf:"concurrency" validate:"required,gt=0"`
				Timeout     time.Duration `koanf:"timeout" validate:"required,gt=0"`
			} `koanf:"handler"`
			Reader struct {
				BatchSize    int64         `koanf:"batch-size" validate:"required,gt=0"`
				BlockTimeout time.Duration `koanf:"block-timeout" validate:"required,gt=0"`
			} `koanf:"reader"`
			Stealer struct {
				Interval           time.Duration `koanf:"interval" validate:"required,gt=0"`
				MinIdleTime        time.Duration `koanf:"min-idle-time" validate:"required,gt=0"`
				MaxDeliveryAttempt int64         `koanf:"max-delivery-attempt" validate:"required,gt=0"`
			} `koanf:"stealer"`
			Reaper struct {
				MinIdleTime time.Duration `koanf:"min-idle-time" validate:"required,gt=0"`
			} `koanf:"reaper"`
		} `koanf:"image-process-request"`

		ImageProcessResult struct {
			StreamKey  string `koanf:"stream-key" validate:"required"`
			StreamSize int    `koanf:"stream-size" validate:"required,gt=0"`
		} `koanf:"image-process-result"`
	} `koanf:"streams"`
}

type AWSConfig struct {
	S3 S3Config `koanf:"s3"`
}

type S3Config struct {
	Bucket string `koanf:"bucket" validate:"required"`
}
