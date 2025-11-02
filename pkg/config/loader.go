package config

import (
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
	"config.default.yaml",
	"config.yaml",
	"config.local.yaml",
}

// Load loads configuration from the given directory into the config struct of
// type T. It loads configuration files in the following order.
//
//  1. config.default.yaml
//  2. config.yaml
//  3. config.local.yaml
//
// At least one of the files must exist. If multiple files exist, later files
// override earlier one's fields.
func Load[T any](dir string) (cfg T, err error) {
	k := koanf.New(".")

	configFiles := lo.Map(
		configFileNames,
		func(f string, _ int) string { return filepath.Join(dir, f) })
	anyFileExists := false
	for _, f := range configFiles {
		_, err := os.Stat(f)
		switch {
		case os.IsNotExist(err):
			continue
		case err != nil:
			return cfg, fmt.Errorf("checking config file %s: %w", f, err)
		}

		anyFileExists = true
		if err := k.Load(file.Provider(f), yaml.Parser()); err != nil {
			return cfg, fmt.Errorf("loading from file(%s): %w", f, err)
		}
	}

	if !anyFileExists {
		return cfg, fmt.Errorf("no config files %v found in %s", configFileNames, dir)
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

	return cfg, nil
}

// LoadValidated loads configuration using [Load] and validates the resulting
// config struct using the validation package.
func LoadValidated[T any](dir string) (cfg T, err error) {
	cfg, err = Load[T](dir)
	if err != nil {
		return cfg, fmt.Errorf("loading config: %w", err)
	}

	if err := validation.Validate(&cfg); err != nil {
		return cfg, fmt.Errorf("validating config struct: %w", err)
	}

	return cfg, nil
}
