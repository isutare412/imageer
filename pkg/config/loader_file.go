package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/samber/lo"
)

var configFileNames = []string{
	"config.default.yaml",
	"config.yaml",
	"config.local.yaml",
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
