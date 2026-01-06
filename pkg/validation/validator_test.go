package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type person struct {
	Address   *address  `validate:"omitempty"`
	Addresses []address `validate:"dive"`
	Names     []string  `validate:"dive,required"`
	Kebab     string    `validate:"omitempty,kebabcase"`
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
				Address:   &address{Street: "foo"},
				Addresses: []address{{Street: "bar"}},
				Names:     []string{"Alice", "Bob"},
			},
			wantErr: false,
		},
		{
			name: "nested struct pointer",
			input: person{
				Address: &address{Street: ""},
			},
			wantErr: true,
		},
		{
			name: "nested struct slice",
			input: person{
				Addresses: []address{{Street: ""}},
			},
			wantErr: true,
		},
		{
			name: "empty slice",
			input: person{
				Addresses: []address{},
			},
			wantErr: false,
		},
		{
			name: "nil slice",
			input: person{
				Addresses: nil,
			},
			wantErr: false,
		},
		{
			name: "blank name",
			input: person{
				Names: []string{""},
			},
			wantErr: true,
		},
		{
			name: "kebab case",
			input: person{
				Kebab: "valid-kebab-case",
			},
			wantErr: false,
		},
		{
			name: "invalid kebab case",
			input: person{
				Kebab: "InvalidKebabCase",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			err := v.Validate(tt.input)
			if tt.wantErr {
				require.Error(t, err, "Validate() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				require.NoError(t, err, "Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
