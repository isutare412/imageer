package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseImageS3Key(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		key           string
		wantProjectID string
		wantImageID   string
		wantOK        bool
	}{
		{
			name:          "normal case",
			key:           "local/foo/projects/5480727a-98a0-4b49-974a-790d2e18e3f5/images/44d2c777-1d83-418c-8359-0a810acaf8cb/original.jpg",
			wantProjectID: "5480727a-98a0-4b49-974a-790d2e18e3f5",
			wantImageID:   "44d2c777-1d83-418c-8359-0a810acaf8cb",
			wantOK:        true,
		},
		{
			name:          "invalid format",
			key:           "local/foo/bar/original.jpg",
			wantProjectID: "",
			wantImageID:   "",
			wantOK:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectID, imageID, ok := parseImageS3Key(tt.key)
			if tt.wantOK {
				require.True(t, ok)
				assert.Equal(t, tt.wantProjectID, projectID)
				assert.Equal(t, tt.wantImageID, imageID)
			} else {
				require.False(t, ok)
			}
		})
	}
}
