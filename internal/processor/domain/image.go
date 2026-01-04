package domain

import (
	"github.com/isutare412/imageer/pkg/images"
)

type RawImage struct {
	Data   []byte
	Format images.Format
}
