package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
)

func loadFromEnv(k *koanf.Koanf, appName string) error {
	// APP_FOO_BAR=baz -> foo.bar=baz
	prefix := appNamePrefixedKey(appName, "")
	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: prefix,
		TransformFunc: func(k, v string) (string, any) {
			k = strings.TrimPrefix(k, prefix)
			k = strings.ToLower(k)
			return strings.ReplaceAll(k, "_", "."), v
		},
	}), nil); err != nil {
		return fmt.Errorf("loading from env: %w", err)
	}

	return nil
}
