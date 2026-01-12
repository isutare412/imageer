package port

import (
	"context"

	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
)

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type ImageProcessRequestQueue interface {
	Push(context.Context, *imageerv1.ImageProcessRequest) error
}
