package valkey

import (
	"strconv"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/pkg/dbhelpers/valkeystream"
)

type ClientConfig struct {
	Addresses []string
	Username  string
	Password  string
}

func (c ClientConfig) applyToOption(opt *valkey.ClientOption) {
	opt.InitAddress = c.Addresses
	opt.Username = c.Username
	opt.Password = c.Password
}

type ImageProcessResultQueueConfig struct {
	StreamKey  string
	StreamSize int
}

func (c ImageProcessResultQueueConfig) StreamSizeString() string {
	return strconv.Itoa(c.StreamSize)
}

type ImageProcessRequestHandlerConfig struct {
	StreamKey            string
	GroupName            string
	HandleConcurrency    int
	HandleTimeout        time.Duration
	ReadBlockTimeout     time.Duration
	ReadBatchSize        int64
	ReapConsumerIdleTime time.Duration
	StealInterval        time.Duration
	StealMinIdleTime     time.Duration
	MaxDeliveryAttempt   int64
}

func (c ImageProcessRequestHandlerConfig) ToConsumerConfig(consumerName string,
) valkeystream.ConsumerConfig {
	return valkeystream.ConsumerConfig{
		Stream: c.StreamKey,
		Group:  c.GroupName,
		Name:   consumerName,
	}
}

func (c ImageProcessRequestHandlerConfig) ToInitializerConfig(consumerName string,
) valkeystream.InitializerConfig {
	return valkeystream.InitializerConfig{
		Consumer: c.ToConsumerConfig(consumerName),
	}
}

func (c ImageProcessRequestHandlerConfig) ToReaperConfig() valkeystream.ReaperConfig {
	return valkeystream.ReaperConfig{
		Stream:            c.StreamKey,
		Group:             c.GroupName,
		IdleTimeThreshold: c.ReapConsumerIdleTime,
	}
}

func (c ImageProcessRequestHandlerConfig) ToReaderConfig(consumerName string,
) valkeystream.ReaderConfig {
	return valkeystream.ReaderConfig{
		Consumer:         c.ToConsumerConfig(consumerName),
		EntryFieldKey:    "msg",
		ReadBlockTimeout: c.ReadBlockTimeout,
		ReadBatchSize:    c.ReadBatchSize,
	}
}

func (c ImageProcessRequestHandlerConfig) ToStealerConfig(consumerName string,
) valkeystream.StealerConfig {
	return valkeystream.StealerConfig{
		Consumer:           c.ToConsumerConfig(consumerName),
		EntryFieldKey:      "msg",
		StealInterval:      c.StealInterval,
		StealMinIdleTime:   c.StealMinIdleTime,
		MaxDeliveryAttempt: c.MaxDeliveryAttempt,
	}
}
