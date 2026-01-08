package image

import (
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

func findPresetNameDiference(requested []string, existing []domain.Preset) []string {
	existingNames := lo.Map(existing, func(p domain.Preset, _ int) string {
		return p.Name
	})

	left, _ := lo.Difference(requested, existingNames)
	return left
}
