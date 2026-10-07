package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSASLMechanism_Validate(t *testing.T) {
	tests := []struct {
		name      string
		mechanism SASLMechanism
		expectErr bool
	}{
		{name: "plain", mechanism: SASLMechanismPlain, expectErr: false},
		{name: "scram-sha-512", mechanism: SASLMechanismScramSHA512, expectErr: false},
		{name: "empty", mechanism: "", expectErr: true},
		{name: "upper case", mechanism: "SCRAM-SHA-512", expectErr: true},
		{name: "unsupported scram variant", mechanism: "scram-sha-256", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mechanism.Validate()

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSASLMechanism_Mechanism(t *testing.T) {
	tests := []struct {
		name      string
		mechanism SASLMechanism
		wantName  string
	}{
		{name: "empty defaults to plain", mechanism: "", wantName: "PLAIN"},
		{name: "plain", mechanism: SASLMechanismPlain, wantName: "PLAIN"},
		{name: "scram-sha-512", mechanism: SASLMechanismScramSHA512, wantName: "SCRAM-SHA-512"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mechanism := tt.mechanism.Mechanism("u", "p")

			assert.Equal(t, tt.wantName, mechanism.Name())
		})
	}
}
