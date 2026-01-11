package csmgroup

import "time"

type ReaderConfig struct {
	Stream           string
	Group            string
	Consumer         string
	EntryFieldKey    string
	ReadBlockTimeout time.Duration
	ReadBatchSize    int64
}

type StealerConfig struct {
	Stream             string
	Group              string
	Consumer           string
	EntryFieldKey      string
	StealInterval      time.Duration
	StealMinIdleTime   time.Duration
	MaxDeliveryAttempt int64
}
