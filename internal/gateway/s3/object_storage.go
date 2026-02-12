package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/trace"

	"github.com/isutare412/imageer/pkg/awshelpers"
	"github.com/isutare412/imageer/pkg/tracing"
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

func (s *ObjectStorage) DeleteObjects(ctx context.Context, keys []string) error {
	ctx, span := tracing.StartSpan(ctx, "s3.ObjectStorage.DeleteObjects",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.PeerServiceAWSS3))
	defer span.End()

	if len(keys) == 0 {
		return nil
	}

	objects := lo.Map(keys, func(key string, _ int) types.ObjectIdentifier {
		return types.ObjectIdentifier{Key: &key}
	})

	_, err := s.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: &s.cfg.Bucket,
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   new(true),
		},
	})
	if err != nil {
		return awshelpers.WrapS3Error(err, "Failed to delete objects")
	}

	return nil
}
