package apperr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/isutare412/imageer/pkg/apperr"
)

func TestError_Unwrap(t *testing.T) {
	tests := []struct {
		name      string
		originErr error
	}{
		{
			name:      "nil_error",
			originErr: nil,
		},
		{
			name:      "non_nil_error",
			originErr: errors.New("test error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apperr.NewError(apperr.CodeBadRequest).WithError(tt.originErr)
			got := errors.Is(err, tt.originErr)
			assert.Equal(t, tt.originErr != nil, got, "Unwrap should return the original error if it exists")
		})
	}
}

func TestAsError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "nil_error",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "apperr_error",
			err:     apperr.NewError(apperr.CodeBadRequest),
			wantErr: true,
		},
		{
			name:    "other_error",
			err:     errors.New("generic error"),
			wantErr: false,
		},
		{
			name:    "wrapped_error",
			err:     fmt.Errorf("wrapped: %w", apperr.NewError(apperr.CodeBadRequest)),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr, ok := apperr.AsError(tt.err)
			if tt.wantErr {
				assert.True(t, ok, "Expected error to be of type *apperr.Error")
				assert.NotNil(t, appErr, "Expected non-nil *apperr.Error")
			} else {
				assert.False(t, ok, "Expected error not to be of type *apperr.Error")
				assert.Nil(t, appErr, "Expected nil *apperr.Error")
			}
		})
	}
}
