package kafka

import (
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Partitioner string

const (
	PartitionerUniformBytes Partitioner = "uniform-bytes"
	PartitionerRoundRobin   Partitioner = "round-robin"
)

func (p Partitioner) Validate() error {
	switch p {
	case PartitionerUniformBytes:
	case PartitionerRoundRobin:
	default:
		return fmt.Errorf("unexpected partitioner %q", p)
	}
	return nil
}

func (p Partitioner) KgoOpt() kgo.Opt {
	switch p {
	case PartitionerRoundRobin:
		return kgo.RecordPartitioner(kgo.RoundRobinPartitioner())
	case PartitionerUniformBytes:
		fallthrough
	default:
		return kgo.RecordPartitioner(kgo.UniformBytesPartitioner(64*1024, true, true, nil))
	}
}

type ClientConfig struct {
	Addrs         []string
	User          string
	Password      string
	ConsumerGroup string
	Partitioner   Partitioner
	ConsumeTopics []string
}

type ImageProcessRequestQueueConfig struct {
	Topic string
}

type ImageProcessResultHandlerConfig struct {
	RetryTopic      string
	HandleTimeout   time.Duration
	MaxRetryAttempt int
	RetryBaseDelay  time.Duration
}

type ImageS3DeleteRequestQueueConfig struct {
	Topic string
}

type ImageS3DeleteRequestHandlerConfig struct {
	RetryTopic      string
	HandleTimeout   time.Duration
	MaxRetryAttempt int
	RetryBaseDelay  time.Duration
}
