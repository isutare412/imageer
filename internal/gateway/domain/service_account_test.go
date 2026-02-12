package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/isutare412/imageer/pkg/serviceaccounts"
	"github.com/isutare412/imageer/pkg/validation"
)

func TestUpdateServiceAccountRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		req     UpdateServiceAccountRequest
		wantErr bool
	}{
		{
			name: "all fields set",
			req: UpdateServiceAccountRequest{
				ID:          "test-id",
				Name:        new("test-name"),
				AccessScope: new(serviceaccounts.AccessScopeProject),
				ProjectIDs:  []string{"project-1", "project-2"},
				ExpireAt:    new(time.Now().Add(time.Hour)),
			},
			wantErr: false,
		},
		{
			name: "empty id",
			req: UpdateServiceAccountRequest{
				ID: "",
			},
			wantErr: true,
		},
		{
			name: "expireAt set to past",
			req: UpdateServiceAccountRequest{
				ID:       "test-id",
				ExpireAt: new(time.Now().Add(-time.Hour)),
			},
			wantErr: true,
		},
		{
			name: "invalid access scope",
			req: UpdateServiceAccountRequest{
				ID:          "test-id",
				AccessScope: new(serviceaccounts.AccessScope("INVALID")),
			},
			wantErr: true,
		},
		{
			name: "blank project IDs",
			req: UpdateServiceAccountRequest{
				ID:         "test-id",
				ProjectIDs: []string{"", ""},
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
