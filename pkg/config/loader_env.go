package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
)

func loadFromEnv(k *koanf.Koanf, appName string) error {
	// APP__FOO__BAR=baz -> foo.bar=baz   (double underscore: nesting level)
	// APP__FOO_BAR=baz  -> foo-bar=baz   (single underscore: word separator)
	prefix := appNamePrefixedKey(appName, "")
	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: prefix,
		TransformFunc: func(k, v string) (string, any) {
			k = strings.TrimPrefix(k, prefix)
			k = strings.ToLower(k)
			k = strings.ReplaceAll(k, "__", ".")
			k = strings.ReplaceAll(k, "_", "-")
			return k, v
		},
	}), nil); err != nil {
		return fmt.Errorf("loading from env: %w", err)
	}

	return nil
}
