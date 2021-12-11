package processor

import "context"

type MsgQueue interface {
	Init(ctx context.Context, topic string) error
	Read(ctx context.Context, topic string, limit int64) (<-chan []byte, error)
}
