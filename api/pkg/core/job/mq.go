package job

import "context"

type MsgQueue interface {
	Produce(ctx context.Context, topic string, val []byte) error
}
