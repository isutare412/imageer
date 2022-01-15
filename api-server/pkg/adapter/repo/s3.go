package repo

import (
	"context"
	"fmt"
	"io"

	"github.com/isutare412/imageer/api-server/pkg/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3 struct {
	client *minio.Client
}

func (s *S3) Put(ctx context.Context, bucket, path string, body io.Reader) error {
	_, err := s.client.PutObject(ctx, bucket, path, body, -1, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("on s3 put: %w", err)
	}
	return nil
}

func (s *S3) Get(ctx context.Context, bucket, path string) (io.ReadSeekCloser, error) {
	obj, err := s.client.GetObject(ctx, bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("on s3 get: %w", err)
	}
	return obj, nil
}

func NewS3(cfg *config.S3Config) (*S3, error) {
	client, err := minio.New(
		cfg.Address,
		&minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: false,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("on new S3 repo: %w", err)
	}

	return &S3{
		client: client,
	}, nil
}
