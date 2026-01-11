package port

import (
	"context"

	"github.com/isutare412/imageer/internal/processor/domain"
)

//go:generate sh -c "go tool mockgen -package $GOPACKAGE -source=$GOFILE -destination=$(basename $GOFILE .go)_mock.go"

type ImageProcessor interface {
	Process(context.Context, domain.RawImage, domain.Preset) (domain.RawImage, error)
}
