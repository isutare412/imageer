package domain

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/isutare412/imageer/pkg/validation"
)

func TestUpsertPresetRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		req     UpsertPresetRequest
		wantErr bool
	}{
		{
			name: "create request with all fields set",
			req: UpsertPresetRequest{
				Name:    lo.ToPtr("w100h100"),
				Default: lo.ToPtr(true),
				Width:   lo.ToPtr[int64](100),
				Height:  lo.ToPtr[int64](100),
			},
			wantErr: false,
		},
		{
			name: "invalid preset dimensions",
			req: UpsertPresetRequest{
				Name:    lo.ToPtr("w100h100"),
				Default: lo.ToPtr(true),
				Width:   lo.ToPtr[int64](0), // invalid width
				Height:  lo.ToPtr[int64](100),
			},
			wantErr: true,
		},
		{
			name: "update request",
			req: UpsertPresetRequest{
				Name:    lo.ToPtr("w100h100"),
				Default: lo.ToPtr(true),
			},
			wantErr: false,
		},
		{
			name: "invalid name format",
			req: UpsertPresetRequest{
				ID:   lo.ToPtr("preset-1"),
				Name: lo.ToPtr("W100H100"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.Validate(tt.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
