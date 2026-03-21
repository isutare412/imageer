package config

import (
	"strings"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/processor/kafka"
	"github.com/isutare412/imageer/internal/processor/s3"
	"github.com/isutare412/imageer/internal/processor/web"
	"github.com/isutare412/imageer/pkg/log"
	"github.com/isutare412/imageer/pkg/tracing"
)

func (c *Config) ToLogConfig() log.Config {
	return log.Config{
		Format:    c.Log.Format,
		Level:     c.Log.Level,
		AddSource: c.Log.AddSource,
		Component: c.Log.Component,
	}
}

func (c *Config) ToTracingConfig() tracing.Config {
	return tracing.Config(c.Trace)
}

func (c *Config) ToS3ObjectStorageConfig() s3.ObjectStorageConfig {
	return s3.ObjectStorageConfig{
		Bucket: c.AWS.S3.Bucket,
	}
}

func (c *Config) ToKafkaClientConfig() kafka.ClientConfig {
	return kafka.ClientConfig{
		Addrs:         parseCSV(c.Kafka.Addresses, ","),
		User:          c.Kafka.Username,
		Password:      c.Kafka.Password,
		ConsumerGroup: c.Kafka.ConsumerGroup,
		ConsumeTopics: []string{
			c.Kafka.Topics.ImageProcessRequest.Topic,
			c.Kafka.Topics.ImageProcessRequest.RetryTopic,
		},
	}
}

func (c *Config) ToKafkaImageProcessRequestHandlerConfig() kafka.ImageProcessRequestHandlerConfig {
	return kafka.ImageProcessRequestHandlerConfig{
		RetryTopic:      c.Kafka.Topics.ImageProcessRequest.RetryTopic,
		HandleTimeout:   c.Kafka.Topics.ImageProcessRequest.Handler.Timeout,
		MaxRetryAttempt: c.Kafka.Topics.ImageProcessRequest.Handler.MaxRetryAttempt,
		RetryBaseDelay:  c.Kafka.Topics.ImageProcessRequest.Handler.RetryBaseDelay,
	}
}

func (c *Config) ToKafkaImageProcessResultQueueConfig() kafka.ImageProcessResultQueueConfig {
	return kafka.ImageProcessResultQueueConfig{
		Topic: c.Kafka.Topics.ImageProcessResult.Topic,
	}
}

func (c *Config) ToWebServerConfig() web.Config {
	return web.Config{
		Port: c.Web.Port,
	}
}

func parseCSV(s string, delim string) []string {
	parts := strings.Split(s, delim)
	parts = lo.Map(parts, func(item string, _ int) string { return strings.TrimSpace(item) })
	parts = lo.Filter(parts, func(item string, _ int) bool {
		return item != ""
	})
	return parts
}
