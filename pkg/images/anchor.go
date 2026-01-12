package images

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/samber/lo"

	"github.com/isutare412/imageer/pkg/apperr"
	imageerv1 "github.com/isutare412/imageer/pkg/protogen/imageer/v1"
)

type Anchor string

const (
	AnchorSmart  Anchor = "SMART"
	AnchorCenter Anchor = "CENTER"
	AnchorNorth  Anchor = "NORTH"
	AnchorSouth  Anchor = "SOUTH"
	AnchorEast   Anchor = "EAST"
	AnchorWest   Anchor = "WEST"
)

// Ensure interfaces are implemented
var (
	_ driver.Valuer = Anchor("")
	_ sql.Scanner   = (*Anchor)(nil)
)

func NewAnchorFromProto(a imageerv1.ImageAnchor) Anchor {
	switch a {
	case imageerv1.ImageAnchor_IMAGE_ANCHOR_SMART:
		return AnchorSmart
	case imageerv1.ImageAnchor_IMAGE_ANCHOR_CENTER:
		return AnchorCenter
	case imageerv1.ImageAnchor_IMAGE_ANCHOR_NORTH:
		return AnchorNorth
	case imageerv1.ImageAnchor_IMAGE_ANCHOR_SOUTH:
		return AnchorSouth
	case imageerv1.ImageAnchor_IMAGE_ANCHOR_EAST:
		return AnchorEast
	case imageerv1.ImageAnchor_IMAGE_ANCHOR_WEST:
		return AnchorWest
	default:
		return ""
	}
}

func (a *Anchor) GetOrDefault() Anchor {
	return lo.FromPtrOr(a, AnchorSmart)
}

func (a Anchor) Validate() error {
	switch a {
	case AnchorSmart:
	case AnchorCenter:
	case AnchorNorth:
	case AnchorSouth:
	case AnchorEast:
	case AnchorWest:
	default:
		return apperr.NewError(apperr.CodeBadRequest).WithSummary("Unexpected image anchor %q", a)
	}
	return nil
}

func (a Anchor) ToProto() imageerv1.ImageAnchor {
	switch a {
	case AnchorSmart:
		return imageerv1.ImageAnchor_IMAGE_ANCHOR_SMART
	case AnchorCenter:
		return imageerv1.ImageAnchor_IMAGE_ANCHOR_CENTER
	case AnchorNorth:
		return imageerv1.ImageAnchor_IMAGE_ANCHOR_NORTH
	case AnchorSouth:
		return imageerv1.ImageAnchor_IMAGE_ANCHOR_SOUTH
	case AnchorEast:
		return imageerv1.ImageAnchor_IMAGE_ANCHOR_EAST
	case AnchorWest:
		return imageerv1.ImageAnchor_IMAGE_ANCHOR_WEST
	default:
		return imageerv1.ImageAnchor_IMAGE_ANCHOR_UNSPECIFIED
	}
}

func (a Anchor) Value() (driver.Value, error) {
	return string(a), nil
}

func (a *Anchor) Scan(value any) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return apperr.NewError(apperr.CodeInternalServerError).
			WithSummary("Invalid value of image anchor: %[1]T(%[1]v)", value)
	}

	*a = Anchor(str)
	return nil
}
