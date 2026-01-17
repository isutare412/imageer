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

type ImageProcessRequestQueueConfig struct {
	StreamKey  string
	StreamSize int
}

func (c ImageProcessRequestQueueConfig) StreamSizeString() string {
	return strconv.Itoa(c.StreamSize)
}

type ImageProcessResultHandlerConfig struct {
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

func (c ImageProcessResultHandlerConfig) ToConsumerConfig(consumerName string,
) valkeystream.ConsumerConfig {
	return valkeystream.ConsumerConfig{
		Stream: c.StreamKey,
		Group:  c.GroupName,
		Name:   consumerName,
	}
}

func (c ImageProcessResultHandlerConfig) ToInitializerConfig(consumerName string,
) valkeystream.InitializerConfig {
	return valkeystream.InitializerConfig{
		Consumer: c.ToConsumerConfig(consumerName),
	}
}

func (c ImageProcessResultHandlerConfig) ToReaperConfig() valkeystream.ReaperConfig {
	return valkeystream.ReaperConfig{
		Stream:            c.StreamKey,
		Group:             c.GroupName,
		IdleTimeThreshold: c.ReapConsumerIdleTime,
	}
}

func (c ImageProcessResultHandlerConfig) ToReaderConfig(consumerName string,
) valkeystream.ReaderConfig {
	return valkeystream.ReaderConfig{
		Consumer:         c.ToConsumerConfig(consumerName),
		EntryFieldKey:    "msg",
		ReadBlockTimeout: c.ReadBlockTimeout,
		ReadBatchSize:    c.ReadBatchSize,
	}
}

func (c ImageProcessResultHandlerConfig) ToStealerConfig(consumerName string,
) valkeystream.StealerConfig {
	return valkeystream.StealerConfig{
		Consumer:           c.ToConsumerConfig(consumerName),
		EntryFieldKey:      "msg",
		StealInterval:      c.StealInterval,
		StealMinIdleTime:   c.StealMinIdleTime,
		MaxDeliveryAttempt: c.MaxDeliveryAttempt,
	}
}

type ImageUploadDoneSubscriberConfig struct {
	ChannelPrefix string
	MaxRetries    int
}

type ImageProcessDoneSubscriberConfig struct {
	ChannelPrefix string
	MaxRetries    int
}

type ImageNotificationPublisherConfig struct {
	UploadDoneChannelPrefix  string
	ProcessDoneChannelPrefix string
}
