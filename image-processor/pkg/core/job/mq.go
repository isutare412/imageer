package job

import "context"

type MsgQueue interface {
	Init(ctx context.Context, topic string) error
	Consume(ctx context.Context, topic string, limit int64) (<-chan []byte, error)
}
