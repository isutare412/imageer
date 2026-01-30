package config

import (
	"strings"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/processor/s3"
	"github.com/isutare412/imageer/internal/processor/valkey"
	"github.com/isutare412/imageer/internal/processor/web"
	"github.com/isutare412/imageer/pkg/log"
)

func (c *Config) ToLogConfig() log.Config {
	return log.Config{
		Format:    c.Log.Format,
		Level:     c.Log.Level,
		AddSource: c.Log.AddSource,
		Component: c.Log.Component,
	}
}

func (c *Config) ToS3ObjectStorageConfig() s3.ObjectStorageConfig {
	return s3.ObjectStorageConfig{
		Bucket: c.AWS.S3.Bucket,
	}
}

func (c *Config) ToValkeyClientConfig() valkey.ClientConfig {
	return valkey.ClientConfig{
		Addresses: parseCSV(c.Valkey.Addresses, ","),
		Username:  c.Valkey.Username,
		Password:  c.Valkey.Password,
	}
}

func (c *Config) ToValkeyImageProcessResultQueueConfig() valkey.ImageProcessResultQueueConfig {
	return valkey.ImageProcessResultQueueConfig{
		StreamKey:  c.Valkey.Streams.ImageProcessResult.StreamKey,
		StreamSize: c.Valkey.Streams.ImageProcessResult.StreamSize,
	}
}

func (c *Config) ToValkeyImageProcessRequestHandlerConfig() valkey.ImageProcessRequestHandlerConfig {
	return valkey.ImageProcessRequestHandlerConfig{
		StreamKey:            c.Valkey.Streams.ImageProcessRequest.StreamKey,
		GroupName:            c.Valkey.Streams.ImageProcessRequest.GroupName,
		HandleConcurrency:    c.Valkey.Streams.ImageProcessRequest.Handler.Concurrency,
		HandleTimeout:        c.Valkey.Streams.ImageProcessRequest.Handler.Timeout,
		ReadBlockTimeout:     c.Valkey.Streams.ImageProcessRequest.Reader.BlockTimeout,
		ReadBatchSize:        c.Valkey.Streams.ImageProcessRequest.Reader.BatchSize,
		ReapConsumerIdleTime: c.Valkey.Streams.ImageProcessRequest.Reaper.MinIdleTime,
		StealInterval:        c.Valkey.Streams.ImageProcessRequest.Stealer.Interval,
		StealMinIdleTime:     c.Valkey.Streams.ImageProcessRequest.Stealer.MinIdleTime,
		MaxDeliveryAttempt:   c.Valkey.Streams.ImageProcessRequest.Stealer.MaxDeliveryAttempt,
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
