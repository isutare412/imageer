package kafka

import (
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sasl/scram"
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

type SASLMechanism string

const (
	SASLMechanismPlain       SASLMechanism = "plain"
	SASLMechanismScramSHA512 SASLMechanism = "scram-sha-512"
)

// Validate accepts an empty value, which means PLAIN.
func (m SASLMechanism) Validate() error {
	switch m {
	case "", SASLMechanismPlain:
	case SASLMechanismScramSHA512:
	default:
		return fmt.Errorf("unexpected sasl mechanism %q", m)
	}
	return nil
}

// Mechanism returns the SASL mechanism for the credentials. Empty value
// falls back to PLAIN.
func (m SASLMechanism) Mechanism(user, password string) sasl.Mechanism {
	switch m {
	case SASLMechanismScramSHA512:
		return scram.Auth{User: user, Pass: password}.AsSha512Mechanism()
	case SASLMechanismPlain:
		fallthrough
	default:
		return plain.Auth{User: user, Pass: password}.AsMechanism()
	}
}

type ClientConfig struct {
	Addrs         []string
	User          string
	Password      string
	SASLMechanism SASLMechanism
	ConsumerGroup string
	Partitioner   Partitioner
	ConsumeTopics []string
}

type ImageProcessRequestHandlerConfig struct {
	RetryTopic      string
	HandleTimeout   time.Duration
	MaxRetryAttempt int
	RetryBaseDelay  time.Duration
}

type ImageProcessResultQueueConfig struct {
	Topic string
}
