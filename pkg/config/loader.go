package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/samber/lo"

	"github.com/isutare412/imageer/pkg/validation"
)

var configFileNames = []string{
	"config.yaml",
	"config.local.yaml",
}

func LoadValidated[T any](dir string) (cfg T, err error) {
	k := koanf.New(".")

	configFiles := lo.Map(
		configFileNames,
		func(f string, _ int) string { return filepath.Join(dir, f) })
	for i, f := range configFiles {
		if i != 0 { // file is optional except first one
			if _, err := os.Stat(f); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					continue
				}
				return cfg, fmt.Errorf("checking config file existence: %w", err)
			}
		}

		if err := k.Load(file.Provider(f), yaml.Parser()); err != nil {
			return cfg, fmt.Errorf("loading from file(%s): %w", f, err)
		}
	}

	// APP_FOO_BAR=baz -> foo.bar=baz
	if err := k.Load(env.Provider("APP_", ".", func(s string) string {
		s = strings.TrimPrefix(s, "APP_")
		s = strings.ToLower(s)
		return strings.ReplaceAll(s, "_", ".")
	}), nil); err != nil {
		return cfg, fmt.Errorf("loading from env: %w", err)
	}

	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return cfg, fmt.Errorf("unmarshaling into config struct: %w", err)
	}

	validator := validation.NewValidator()
	if err := validator.Validate(&cfg); err != nil {
		return cfg, fmt.Errorf("validating config struct: %w", err)
	}

	return cfg, nil
}
