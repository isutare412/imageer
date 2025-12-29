package apikey

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPIKey_Validate(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		keyGen  func() APIKey
		key     string
		wantErr bool
	}{
		{
			name:    "valid key",
			key:     "ak_SOEXTZL3Ww7BXyqMWdD5QpPfyPF5nXTy14PkpV9X6vySOzmQdpjHYXK26vbhDgDj",
			wantErr: false,
		},
		{
			name:    "invalid length",
			key:     "ak_SOEXTZL3Ww7BXyqMWdD5QpPfyPF5nXTy14PkpV9X6vySOzmQdpjHYXK26vbhDgDjx",
			wantErr: true,
		},
		{
			name:    "invalid prefix",
			key:     "xx_SOEXTZL3Ww7BXyqMWdD5QpPfyPF5nXTy14PkpV9X6vySOzmQdpjHYXK26vbhDgDj",
			wantErr: true,
		},
		{
			name:    "invalid checksum",
			key:     "ak_SOEXTZL3Ww7BXyqMWdD5QpPfyPF5nXTy14PkpV9X6vySOzmQdpjHYXK26vbhDgDk",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := APIKey(tt.key)
			err := k.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
