package image

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

func Test_findPresetNameDiference(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		requested []string
		existing  []domain.Preset
		want      []string
	}{
		{
			name: "all exist",
			requested: []string{
				"preset-1",
				"preset-2",
			},
			existing: []domain.Preset{
				{Name: "preset-1"},
				{Name: "preset-2"},
			},
			want: []string{},
		},
		{
			name: "some not exist",
			requested: []string{
				"preset-1",
				"preset-2",
			},
			existing: []domain.Preset{
				{Name: "preset-1"},
				{Name: "preset-3"},
			},
			want: []string{"preset-2"},
		},
		{
			name: "none exist",
			requested: []string{
				"preset-1",
				"preset-2",
			},
			existing: []domain.Preset{},
			want:     []string{"preset-1", "preset-2"},
		},
		{
			name:      "none requested",
			requested: []string{},
			existing: []domain.Preset{
				{Name: "preset-1"},
				{Name: "preset-2"},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findPresetNameDiference(tt.requested, tt.existing)
			assert.Equal(t, tt.want, got)
		})
	}
}
