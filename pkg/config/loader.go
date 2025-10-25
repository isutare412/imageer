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

func LoadValidated[T any](dir string) (cfg T, err error) {
	k := koanf.New(".")

	configFiles := lo.Map(
		configFileNames,
		func(f string, _ int) string { return filepath.Join(dir, f) })
	for i, f := range configFiles {
		_, err := os.Stat(f)
		switch {
		case i == 0 && os.IsNotExist(err):
			return cfg, fmt.Errorf("default config file %s does not exist: %w", f, err)
		case i > 0 && os.IsNotExist(err):
			continue // optional file missing, skip
		case err != nil:
			return cfg, fmt.Errorf("checking config file %s: %w", f, err)
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

	if err := validation.Validate(&cfg); err != nil {
		return cfg, fmt.Errorf("validating config struct: %w", err)
	}

	return cfg, nil
}
