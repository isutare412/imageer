package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSLMode_Validate(t *testing.T) {
	tests := []struct {
		name      string
		mode      SSLMode
		expectErr bool
	}{
		{name: "disable", mode: SSLModeDisable, expectErr: false},
		{name: "allow", mode: SSLModeAllow, expectErr: false},
		{name: "prefer", mode: SSLModePrefer, expectErr: false},
		{name: "require", mode: SSLModeRequire, expectErr: false},
		{name: "empty defaults to disable", mode: "", expectErr: false},
		{name: "verify-full is not supported", mode: "verify-full", expectErr: true},
		{name: "upper case", mode: "REQUIRE", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mode.Validate()

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildDSN(t *testing.T) {
	cfg := ClientConfig{
		Host:     "db.example",
		Port:     5432,
		User:     "tester",
		Password: "secret",
		Database: "imageer",
	}

	t.Run("empty sslmode falls back to disable", func(t *testing.T) {
		assert.Equal(t,
			"host=db.example port=5432 user=tester password=secret dbname=imageer sslmode=disable",
			buildDSN(cfg))
	})

	t.Run("configured sslmode is used", func(t *testing.T) {
		c := cfg
		c.SSLMode = SSLModeRequire
		assert.Equal(t,
			"host=db.example port=5432 user=tester password=secret dbname=imageer sslmode=require",
			buildDSN(c))
	})
}
