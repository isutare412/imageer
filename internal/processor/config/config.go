package config

import (
	"time"

	"github.com/isutare412/imageer/pkg/log"
)

type Config struct {
	Log   LogConfig   `koanf:"log"`
	Trace TraceConfig `koanf:"trace"`
	Web   WebConfig   `koanf:"web"`
	Kafka KafkaConfig `koanf:"kafka"`
	AWS   AWSConfig   `koanf:"aws"`
}

type LogConfig struct {
	Format    log.Format `koanf:"format" validate:"validateFn=Validate"`
	Level     log.Level  `koanf:"level" validate:"validateFn=Validate"`
	AddSource bool       `koanf:"add-source"`
	Component string     `koanf:"component" validate:"required"`
}

type TraceConfig struct {
	Enabled          bool    `koanf:"enabled"`
	ServiceName      string  `koanf:"service-name" validate:"required"`
	SamplingRatio    float64 `koanf:"sampling-ratio" validate:"required,gte=0,lte=1"`
	OTLPGRPCEndpoint string  `koanf:"otlp-grpc-endpoint" validate:"required"`
}

type WebConfig struct {
	Port int `koanf:"port" validate:"required,gt=0,lte=65535"`
}

type KafkaConfig struct {
	Addresses     string `koanf:"addresses" validate:"required"`
	Username      string `koanf:"username" validate:"required"`
	Password      string `koanf:"password" validate:"required"`
	ConsumerGroup string `koanf:"consumer-group" validate:"required"`

	Topics struct {
		ImageProcessRequest struct {
			Topic      string `koanf:"topic" validate:"required"`
			RetryTopic string `koanf:"retry-topic" validate:"required"`
			Handler    struct {
				Timeout         time.Duration `koanf:"timeout" validate:"required,gt=0"`
				MaxRetryAttempt int           `koanf:"max-retry-attempt" validate:"required,gte=0"`
				RetryBaseDelay  time.Duration `koanf:"retry-base-delay" validate:"required,gt=0"`
			} `koanf:"handler"`
		} `koanf:"image-process-request"`

		ImageProcessResult struct {
			Topic string `koanf:"topic" validate:"required"`
		} `koanf:"image-process-result"`
	} `koanf:"topics"`
}

type AWSConfig struct {
	S3 S3Config `koanf:"s3"`
}

type S3Config struct {
	Bucket string `koanf:"bucket" validate:"required"`
}
