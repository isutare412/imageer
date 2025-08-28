package config_test

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isutare412/imageer/pkg/config"
)

//go:embed testdata/*
var testFS embed.FS

type configModel struct {
	Name   string `validate:"required"`
	Nested nestedConfig
}

type nestedConfig struct {
	Age int `validate:"gt=0"`
}

func copyTestFile(t *testing.T, src, dst string) {
	t.Helper()
	content, err := testFS.ReadFile(src)
	require.NoError(t, err)
	err = os.WriteFile(dst, content, 0644)
	require.NoError(t, err)
}

func TestLoadValidated(t *testing.T) {
	tests := []struct {
		name    string
		want    configModel
		wantErr bool
		setup   func(t *testing.T, dir string)
	}{
		{
			name: "valid config with required fields",
			want: configModel{
				Name:   "test-app",
				Nested: nestedConfig{Age: 25},
			},
			wantErr: false,
			setup: func(t *testing.T, dir string) {
				copyTestFile(t, "testdata/valid_config.yaml", filepath.Join(dir, "config.yaml"))
			},
		},
		{
			name: "config with local override",
			want: configModel{
				Name:   "local-app",
				Nested: nestedConfig{Age: 30},
			},
			wantErr: false,
			setup: func(t *testing.T, dir string) {
				copyTestFile(t, "testdata/valid_config.yaml", filepath.Join(dir, "config.yaml"))
				copyTestFile(t, "testdata/valid_local.yaml", filepath.Join(dir, "config.local.yaml"))
			},
		},
		{
			name: "config with environment variables",
			want: configModel{
				Name:   "env-app",
				Nested: nestedConfig{Age: 35},
			},
			wantErr: false,
			setup: func(t *testing.T, dir string) {
				copyTestFile(t, "testdata/valid_config.yaml", filepath.Join(dir, "config.yaml"))
				t.Setenv("APP_NAME", "env-app")
				t.Setenv("APP_NESTED_AGE", "35")
			},
		},
		{
			name:    "missing required field",
			want:    configModel{},
			wantErr: true,
			setup: func(t *testing.T, dir string) {
				copyTestFile(t, "testdata/missing_required.yaml", filepath.Join(dir, "config.yaml"))
			},
		},
		{
			name:    "invalid validation constraint",
			want:    configModel{},
			wantErr: true,
			setup: func(t *testing.T, dir string) {
				copyTestFile(t, "testdata/invalid_validation.yaml", filepath.Join(dir, "config.yaml"))
			},
		},
		{
			name:    "missing config directory",
			want:    configModel{},
			wantErr: true,
			setup:   func(t *testing.T, dir string) {},
		},
		{
			name:    "invalid yaml syntax",
			want:    configModel{},
			wantErr: true,
			setup: func(t *testing.T, dir string) {
				copyTestFile(t, "testdata/invalid_yaml.yaml", filepath.Join(dir, "config.yaml"))
			},
		},
		{
			name: "only local config file exists",
			want: configModel{
				Name:   "local-app",
				Nested: nestedConfig{Age: 30},
			},
			wantErr: true,
			setup: func(t *testing.T, dir string) {
				copyTestFile(t, "testdata/valid_local.yaml", filepath.Join(dir, "config.local.yaml"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()
			if tt.setup != nil {
				tt.setup(t, testDir)
			}

			got, err := config.LoadValidated[configModel](testDir)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
