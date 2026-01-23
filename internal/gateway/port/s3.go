package port

import (
	"context"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type S3Presigner interface {
	PresignPutObject(context.Context, domain.PresignPutObjectRequest) (domain.PresignPutObjectResponse, error)
}

type ObjectStorage interface {
	DeleteObjects(ctx context.Context, keys []string) error
}
