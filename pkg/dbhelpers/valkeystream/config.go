package valkeystream

import "time"

type ConsumerConfig struct {
	Stream string
	Group  string
	Name   string
}

type InitializerConfig struct {
	Consumer ConsumerConfig
}

type ReaperConfig struct {
	Stream            string
	Group             string
	IdleTimeThreshold time.Duration
}

type ReaderConfig struct {
	Consumer         ConsumerConfig
	EntryFieldKey    string
	ReadBlockTimeout time.Duration
	ReadBatchSize    int64
}

type StealerConfig struct {
	Consumer           ConsumerConfig
	EntryFieldKey      string
	StealInterval      time.Duration
	StealMinIdleTime   time.Duration
	MaxDeliveryAttempt int64
}
