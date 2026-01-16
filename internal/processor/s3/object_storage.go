package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/samber/lo"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/awshelpers"
)

type ObjectStorage struct {
	client *s3.Client
	cfg    ObjectStorageConfig
}

func NewObjectStorage(cfg ObjectStorageConfig) (*ObjectStorage, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("loading default aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)

	return &ObjectStorage{
		client: client,
		cfg:    cfg,
	}, nil
}

func (s *ObjectStorage) Get(ctx context.Context, key string) ([]byte, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.cfg.Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, awshelpers.WrapS3Error(err, "Failed to get object %s", key)
	}
	defer func() { _ = output.Body.Close() }()

	bodyBytes, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, apperr.NewError(apperr.CodeInternalServerError).WithCause(err)
	}

	return bodyBytes, nil
}

func (s *ObjectStorage) Put(ctx context.Context, key string, data []byte,
	contentType string,
) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.cfg.Bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: lo.EmptyableToPtr(contentType),
	})
	if err != nil {
		return awshelpers.WrapS3Error(err, "Failed to put object %s", key)
	}

	return nil
}
