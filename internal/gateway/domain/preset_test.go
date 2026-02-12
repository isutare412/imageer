package domain

import (
	"testing"

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
				Name:    new("w100h100"),
				Default: new(true),
				Width:   new(int64(100)),
				Height:  new(int64(100)),
			},
			wantErr: false,
		},
		{
			name: "invalid preset dimensions",
			req: UpsertPresetRequest{
				Name:    new("w100h100"),
				Default: new(true),
				Width:   new(int64(0)), // invalid width
				Height:  new(int64(100)),
			},
			wantErr: true,
		},
		{
			name: "update request",
			req: UpsertPresetRequest{
				Name:    new("w100h100"),
				Default: new(true),
			},
			wantErr: false,
		},
		{
			name: "invalid name format",
			req: UpsertPresetRequest{
				ID:   new("preset-1"),
				Name: new("W100H100"),
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
