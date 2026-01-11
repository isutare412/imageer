package valkey

import (
	"strconv"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/isutare412/imageer/internal/gateway/valkey/csmgroup"
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

type ImageEventQueueConfig struct {
	StreamKey  string
	StreamSize int
}

func (c ImageEventQueueConfig) StreamSizeString() string {
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

func (c ImageProcessResultHandlerConfig) ToReaderConfig(consumerName string) csmgroup.ReaderConfig {
	return csmgroup.ReaderConfig{
		Stream:           c.StreamKey,
		Group:            c.GroupName,
		Consumer:         consumerName,
		EntryFieldKey:    "msg",
		ReadBlockTimeout: c.ReadBlockTimeout,
		ReadBatchSize:    c.ReadBatchSize,
	}
}

func (c ImageProcessResultHandlerConfig) ToStealerConfig(consumerName string) csmgroup.StealerConfig {
	return csmgroup.StealerConfig{
		Stream:             c.StreamKey,
		Group:              c.GroupName,
		Consumer:           consumerName,
		EntryFieldKey:      "msg",
		StealInterval:      c.StealInterval,
		StealMinIdleTime:   c.StealMinIdleTime,
		MaxDeliveryAttempt: c.MaxDeliveryAttempt,
	}
}
