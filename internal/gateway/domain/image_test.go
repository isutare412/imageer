package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/isutare412/imageer/pkg/images"
	"github.com/isutare412/imageer/pkg/validation"
)

func TestCreatePresignedURLRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		req     CreatePresignedURLRequest
		wantErr bool
	}{
		{
			name: "normal case",
			req: CreatePresignedURLRequest{
				FileName:            "example.jpg",
				ContentType:         images.ContentTypeImageJPEG,
				TransformationNames: []string{"w100h100", "w200h200"},
			},
			wantErr: false,
		},
		{
			name: "invalid content type",
			req: CreatePresignedURLRequest{
				FileName:            "example.jpg",
				ContentType:         images.ContentType("invalid"), // invalid content type
				TransformationNames: []string{"w100h100", "w200h200"},
			},
			wantErr: true,
		},
		{
			name: "empty transformation name",
			req: CreatePresignedURLRequest{
				FileName:            "example.jpg",
				ContentType:         images.ContentTypeImageJPEG,
				TransformationNames: []string{"", "w200h200"}, // invalid transformation name
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
