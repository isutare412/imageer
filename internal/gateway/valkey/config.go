package valkey

import (
	"strconv"

	"github.com/valkey-io/valkey-go"
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
