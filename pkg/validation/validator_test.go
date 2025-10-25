package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type person struct {
	Address   address
	Addresses []address `validate:"dive"`
}

type address struct {
	Street string `validate:"required"`
}

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name: "all set",
			input: person{
				Address:   address{Street: "foo"},
				Addresses: []address{{Street: "bar"}},
			},
			wantErr: false,
		},
		{
			name: "nested struct",
			input: person{
				Address: address{Street: ""},
			},
			wantErr: true,
		},
		{
			name: "nested struct slice",
			input: person{
				Address:   address{Street: "foo"},
				Addresses: []address{{Street: ""}},
			},
			wantErr: true,
		},
		{
			name: "empty slice",
			input: person{
				Address:   address{Street: "foo"},
				Addresses: []address{},
			},
			wantErr: false,
		},
		{
			name: "nil slice",
			input: person{
				Address:   address{Street: "foo"},
				Addresses: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			gotErr := v.Validate(tt.input)
			if tt.wantErr {
				require.Error(t, gotErr, "Validate() error = %v, wantErr %v", gotErr, tt.wantErr)
			} else {
				require.NoError(t, gotErr, "Validate() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}
