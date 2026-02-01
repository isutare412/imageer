package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/pkg/awshelpers"
	"github.com/isutare412/imageer/pkg/trace"
)

type Presigner struct {
	client *s3.PresignClient
	cfg    PresignerConfig
}

func NewPresigner(cfg PresignerConfig) (*Presigner, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("loading default aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)
	presigner := s3.NewPresignClient(client)

	return &Presigner{
		client: presigner,
		cfg:    cfg,
	}, nil
}

func (p *Presigner) PresignPutObject(ctx context.Context, req domain.PresignPutObjectRequest,
) (domain.PresignPutObjectResponse, error) {
	ctx, span := trace.StartSpan(ctx, "s3.Presigner.PresignPutObject")
	defer span.End()

	resp, err := p.client.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      &p.cfg.Bucket,
		Key:         &req.S3Key,
		ContentType: lo.EmptyableToPtr(req.ContentType),
	}, s3.WithPresignExpires(p.cfg.Expiry))
	if err != nil {
		return domain.PresignPutObjectResponse{}, awshelpers.WrapS3Error(err, "Failed to presign")
	}

	// Ensure Content-Type is included in the header for the client to send
	header := resp.SignedHeader
	if req.ContentType != "" {
		header.Set("Content-Type", req.ContentType)
	}

	return domain.PresignPutObjectResponse{
		URL:      resp.URL,
		Header:   header,
		ExpireAt: time.Now().UTC().Add(p.cfg.Expiry - 5*time.Second), // Subtract 5 seconds as buffer
	}, nil
}
