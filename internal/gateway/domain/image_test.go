package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/isutare412/imageer/pkg/images"
	"github.com/isutare412/imageer/pkg/validation"
)

func TestCreateUploadURLRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		req     CreateUploadURLRequest
		wantErr bool
	}{
		{
			name: "normal case",
			req: CreateUploadURLRequest{
				FileName:    "example.jpg",
				Format:      images.FormatWebp,
				PresetNames: []string{"w100h100", "w200h200"},
			},
			wantErr: false,
		},
		{
			name: "invalid content type",
			req: CreateUploadURLRequest{
				FileName:    "example.jpg",
				Format:      images.Format("invalid"), // invalid content type
				PresetNames: []string{"w100h100", "w200h200"},
			},
			wantErr: true,
		},
		{
			name: "empty preset name",
			req: CreateUploadURLRequest{
				FileName:    "example.jpg",
				Format:      images.FormatWebp,
				PresetNames: []string{"", "w200h200"}, // invalid preset name
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
