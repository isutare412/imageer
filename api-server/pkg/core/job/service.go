package job

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/isutare412/imageer/api-server/pkg/config"
)

type Service interface {
	Archive(ctx context.Context, ext string, body io.Reader) error
	Produce(ctx context.Context, val string) error
}

type service struct {
	mq           MsgQueue
	reqQueueName string
	resQueueName string

	objRepo   ObjectRepo
	bucket    string
	sourceDir string
}

func (s *service) Archive(ctx context.Context, ext string, body io.Reader) error {
	if err := s.objRepo.Put(ctx, s.bucket, s.sourceFileName(ext), body); err != nil {
		return fmt.Errorf("on produce bytes: %w", err)
	}
	return nil
}

func (s *service) Produce(ctx context.Context, val string) error {
	if err := s.mq.Produce(ctx, s.reqQueueName, []byte(val)); err != nil {
		return err
	}
	return nil
}

func (s *service) sourceFileName(ext string) string {
	now := time.Now()
	return fmt.Sprintf("%s/%04d/%02d/%02d/%s.%s",
		s.sourceDir, now.Year(), now.Month(), now.Day(), uuid.NewString(), ext)
}

func NewService(cfg *config.JobConfig, mq MsgQueue, objRepo ObjectRepo) Service {
	return &service{
		mq:           mq,
		reqQueueName: cfg.Queue.Request,
		resQueueName: cfg.Queue.Response,
		objRepo:      objRepo,
		bucket:       cfg.Repo.S3.Bucket,
		sourceDir:    cfg.Repo.S3.SourceDir,
	}
}
