package config

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/knadh/koanf/v2"

	"github.com/isutare412/imageer/pkg/validation"
)

// Load loads configuration into a struct of type T from multiple sources.
// Order of precedence (highest to lowest):
//  1. Environment variables (prefixed with {appName}_)
//  2. AWS SSM Parameter Store (if enabled via {appName}_SSM_ENABLED env var)
//  3. Configuration files in the specified dir directory in the following order:
//     config.local.yaml, config.yaml, config.default.yaml
func Load[T any](appName, dir string) (cfg T, err error) {
	k := koanf.New(".")

	if err := loadFromFile(k, dir); err != nil {
		return cfg, fmt.Errorf("loading from files: %w", err)
	}

	if err := loadFromAWSParameterStore(k, appName); err != nil {
		return cfg, fmt.Errorf("loading from AWS SSM parameter store: %w", err)
	}

	if err := loadFromEnv(k, appName); err != nil {
		return cfg, fmt.Errorf("loading from env: %w", err)
	}

	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return cfg, fmt.Errorf("unmarshaling into config struct: %w", err)
	}

	return cfg, nil
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

func appNamePrefixedKey(appName, key string) string {
	appName = cmp.Or(appName, "APP") // default to APP if empty
	appName = strings.ReplaceAll(appName, "-", "_")
	appName = strings.ToUpper(appName)
	return fmt.Sprintf("%s_%s", appName, key)
}
