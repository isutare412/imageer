package domain

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/isutare412/imageer/pkg/validation"
)

func TestUpsertTransformationRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		req     UpsertTransformationRequest
		wantErr bool
	}{
		{
			name: "create request with all fields set",
			req: UpsertTransformationRequest{
				Name:    lo.ToPtr("w100h100"),
				Default: lo.ToPtr(true),
				Width:   lo.ToPtr[int64](100),
				Height:  lo.ToPtr[int64](100),
			},
			wantErr: false,
		},
		{
			name: "invalid transformation dimensions",
			req: UpsertTransformationRequest{
				Name:    lo.ToPtr("w100h100"),
				Default: lo.ToPtr(true),
				Width:   lo.ToPtr[int64](0), // invalid width
				Height:  lo.ToPtr[int64](100),
			},
			wantErr: true,
		},
		{
			name: "update request",
			req: UpsertTransformationRequest{
				ID:      lo.ToPtr("w100h100"),
				Default: lo.ToPtr(true),
			},
			wantErr: false,
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
