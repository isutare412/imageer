package job

import (
	"context"
	"io"
)

type ObjectRepo interface {
	Put(ctx context.Context, bucket, path string, body io.Reader) error
	Get(ctx context.Context, bucket, path string) (io.ReadSeekCloser, error)
}
