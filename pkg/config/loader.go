package config

import (
	"cmp"
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

// Load loads configuration into a struct of type T from multiple sources.
// Order of precedence (highest to lowest):
//  1. Environment variables (prefixed with {appName}_)
//  2. Configuration files in the specified dir directory in the following order:
//     config.local.yaml, config.yaml, config.default.yaml
func Load[T any](appName, dir string) (cfg T, err error) {
	k := koanf.New(".")

	if err := loadFromFile(k, dir); err != nil {
		return cfg, fmt.Errorf("loading from files: %w", err)
	}

	if err := loadFromEnv(k, appName); err != nil {
		return cfg, fmt.Errorf("loading from env: %w", err)
	}

	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return cfg, fmt.Errorf("unmarshaling into config struct: %w", err)
	}

	return cfg, nil
}

func loadFromFile(k *koanf.Koanf, dir string) error {
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
			return fmt.Errorf("checking config file %s: %w", f, err)
		}

		anyFileExists = true
		if err := k.Load(file.Provider(f), yaml.Parser()); err != nil {
			return fmt.Errorf("loading from file(%s): %w", f, err)
		}
	}

	if !anyFileExists {
		return fmt.Errorf("no config files %v found in %s", configFileNames, dir)
	}

	return nil
}

func loadFromEnv(k *koanf.Koanf, appName string) error {
	// APP_FOO_BAR=baz -> foo.bar=baz
	prefix := appNamePrefixedKey(appName, "")
	if err := k.Load(env.Provider(prefix, ".", func(s string) string {
		s = strings.TrimPrefix(s, prefix)
		s = strings.ToLower(s)
		return strings.ReplaceAll(s, "_", ".")
	}), nil); err != nil {
		return fmt.Errorf("loading from env: %w", err)
	}

	return nil
}

func appNamePrefixedKey(appName, key string) string {
	appName = cmp.Or(appName, "APP") // default to APP if empty
	appName = strings.ReplaceAll(appName, "-", "_")
	appName = strings.ToUpper(appName)
	return fmt.Sprintf("%s_%s", appName, key)
}

// LoadValidated loads configuration using [Load] and validates the resulting
// config struct using the validation package.
func LoadValidated[T any](appName, dir string) (cfg T, err error) {
	cfg, err = Load[T](appName, dir)
	if err != nil {
		return cfg, fmt.Errorf("loading config: %w", err)
	}

	if err := validation.Validate(&cfg); err != nil {
		return cfg, fmt.Errorf("validating config struct: %w", err)
	}

	return cfg, nil
}
