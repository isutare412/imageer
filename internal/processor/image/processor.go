package image

import (
	"context"
	"fmt"

	"github.com/h2non/bimg"

	"github.com/isutare412/imageer/internal/processor/domain"
	"github.com/isutare412/imageer/pkg/tracing"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (c *Processor) Process(ctx context.Context, input domain.RawImage, preset domain.Preset,
) (domain.RawImage, error) {
	_, span := tracing.StartSpan(ctx, "image.Processor.Process")
	defer span.End()

	var opt bimg.Options
	applyPreset(&opt, preset)

	img := bimg.NewImage(input.Data)
	outBytes, err := img.Process(opt)
	if err != nil {
		return domain.RawImage{}, wrapBimgError(err, "Failed to process image")
	}

	meta, err := img.Metadata()
	if err != nil {
		return domain.RawImage{}, wrapBimgError(err, "Failed to get image metadata")
	}

	format, err := imageTypeToFormat(meta.Type)
	if err != nil {
		return domain.RawImage{}, fmt.Errorf("converting image type to format: %w", err)
	}

	return domain.RawImage{
		Data:   outBytes,
		Format: format,
	}, nil
}
